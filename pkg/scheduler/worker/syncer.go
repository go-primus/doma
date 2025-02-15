package worker

import (
	"time"

	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/go-primus/doma/pkg/scheduler/internal/store"
)

// syncer 用于任务结果和状态同步

type syncer struct {
	sync <-chan *model.SyncMessage

	//
	done   chan struct{}
	notify chan struct{}

	//
	store    store.TaskStore
	syncFunc func(status model.SyncMessage) error
}

type syncerParams struct {
	sync  <-chan *model.SyncMessage
	store store.TaskStore
}

func newSyncer(params syncerParams) *syncer {
	return &syncer{
		done:   make(chan struct{}),
		notify: make(chan struct{}),
		sync:   params.sync,
		store:  params.store,
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
			s.store.UpdateStatus(syncmsg.Task.ID, syncmsg.Status)

			if syncmsg.Status.Status == model.TaskState_Running {
				s.store.UpdateProgress(syncmsg.Task.ID, syncmsg.Progress)
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
		err := s.syncFunc(model.SyncMessage{
			Task:     task.Msg,
			Status:   task.Status,
			Progress: task.Progress,
		})
		if err == nil {
			s.store.MarkSynced(task.Msg.ID)
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
