package scheduler

import (
	"time"
)

// syncer 用于任务结果和状态同步

type syncer struct {
	sync <-chan *SyncMessage

	//
	done chan struct{}

	//
	store    TaskStore
	syncFunc func(status TaskStatus) error
}

func (s *syncer) shutdown() {
	s.done <- struct{}{}
}

func (s *syncer) Start() {

	// go func() {
	// 	s.syncing()
	// }()

	go func() {
		s.syncTask()
	}()
}

func (s *syncer) syncTask() {
	for {
		select {
		case <-s.done:
			// trigger sync
			s.syncinternal()
			return
		case <-s.sync:
		// 	// triger sync

		// 	s.syncinternal()
		case <-time.After(time.Second):
			// trigger sync

			s.syncinternal()
		}
	}
}

func (s *syncer) syncinternal() {
	// 获取未同步状态

	tasks := s.store.ListSyncStatus()
	for _, task := range tasks {
		err := s.syncFunc(task)
		if err == nil {
			s.store.MarkSynced(task.ID)
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
