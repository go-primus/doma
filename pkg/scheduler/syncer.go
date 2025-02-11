package scheduler

import (
	"time"
)

// syncer 用于任务结果和状态同步

type syncer struct {
	sync <-chan *SyncMessage

	//
	done   chan struct{}
	notify chan struct{}

	//
	store    TaskStore
	syncFunc func(status SyncMessage) error
}

func newSyncer() *syncer {
	return &syncer{
		done:   make(chan struct{}),
		notify: make(chan struct{}),
	}
}

func (s *syncer) shutdown() {
	s.done <- struct{}{}
}

func (s *syncer) Start() {

	go func() {
		s.syncing()
	}()

	go func() {
		s.syncTask()
	}()
}

func (s *syncer) syncing() {
	for {
		select {
		case <-s.done:
			return
		case syncmsg := <-s.sync:
			s.store.UpdateStatus(syncmsg.task.ID, syncmsg.status)

			if syncmsg.status.Status == TaskState_Running {
				s.store.UpdateProgress(syncmsg.task.ID, syncmsg.progress)
			}
			s.notify <- struct{}{}

		}
	}
}

func (s *syncer) syncTask() {
	for {
		select {
		case <-s.done:
			// trigger sync
			s.syncinternal()
			return
		case <-s.notify:
			s.syncinternal()
		case <-time.After(time.Second):
			// trigger sync
			s.syncinternal()
		}
	}
}

func (s *syncer) syncinternal() {
	// 获取未同步状态

	tasks := s.store.ListSyncTasks()
	for _, task := range tasks {
		err := s.syncFunc(SyncMessage{
			task:     task.msg,
			status:   task.status,
			progress: task.progress,
		})
		if err == nil {
			s.store.MarkSynced(task.msg.ID)
		}
	}
}

// func (s *syncer) syncing() {
// 	var requests []*syncFunc

// 	for {
// 		select {
// 		case <-s.done:
// 			for _, req := range requests {
// 				if err := s.syncFunc(req.status); err != nil {
// 					slog.Error(err.Error())
// 				}
// 			}
// 			slog.Debug("syncer done")
// 			return
// 		case syncMsg := <-s.sync:
// 			requests = append(requests, &syncFunc{
// 				SyncMessage: *syncMsg,
// 				fn: func() error {
// 					return nil
// 				},
// 			})
// 		case <-time.After(time.Second * 1):
// 			var temp []*syncFunc
// 			for _, req := range requests {
// 				// if req.deadline
// 				if err := s.syncFunc(req.status); err != nil {
// 					temp = append(temp, req)
// 				}
// 			}
// 			requests = temp
// 		}
// 	}
// }

// type syncFunc struct {
// 	SyncMessage
// 	fn func() error
// }
