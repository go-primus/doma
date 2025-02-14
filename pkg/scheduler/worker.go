package scheduler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/go-primus/doma/pkg/scheduler/internal/cancelation"
	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/go-primus/doma/pkg/scheduler/internal/store"
)

type worker struct {
	bus eventbus.EventBus

	dispatch chan<- model.TaskMessage // dispatch queue

	//
	processor  *processor
	syncer     *syncer
	subscriber *subscriber

	//
	cancelCh chan<- model.TaskCommand

	//
	finished chan<- model.TaskMessage

	//
	store store.TaskStore
}

func NewWorker(bus eventbus.EventBus, handler Handler) *worker {

	store := store.NewTaskStore()
	dispatch := make(chan model.TaskMessage, 10)
	syncCh := make(chan *model.SyncMessage)

	cancels := cancelation.NewCancelations()
	// processor
	processor := NewProcessor(dispatch, cancels)
	processor.sync = syncCh
	processor.handler = handler
	processor.store = store

	// syncer
	syncer := newSyncer()
	syncer.sync = syncCh
	syncer.store = store

	syncer.syncFunc = func(syncmsg model.SyncMessage) error {
		slog.Info("sync task", "task", syncmsg.Task.ID, "status", syncmsg.Status.Status, "---", syncmsg.Progress.Progress)
		return nil
	}

	cancelCh := make(chan model.TaskCommand)

	subscriber := newSubscriber(subscriberParams{
		cancel:       cancelCh,
		cancelations: cancels,
	})

	return &worker{
		bus:        bus,
		dispatch:   dispatch,
		processor:  processor,
		store:      store,
		syncer:     syncer,
		subscriber: subscriber,
		cancelCh:   cancelCh,
	}
}

func (s *worker) Start() {
	s.syncer.Start()
	s.processor.Start()
	s.subscriber.Start()
	s.start()

}

func (s *worker) start() {
	s.bus.Subscribe("tasks.queues", func(msg *eventbus.Msg) {
		var taskMsg model.TaskMessage
		json.Unmarshal(msg.Data, &taskMsg)

		fmt.Println("handle task queues msg")

		// 任务去重
		// 判断任务是否在执行，任务已经执行，则忽略
		//
		// status := s.store.GetStatus(taskMsg.ID)
		// check task
		// 查询任务是否正在执行，如果正在执行忽略
		//
		s.store.AddTask(taskMsg)
		s.dispatch <- taskMsg
		fmt.Println("handle task queues msgxxx")
	})

	s.bus.Subscribe("tasks.commands", func(msg *eventbus.Msg) {
		var taskCommand model.TaskCommand
		json.Unmarshal(msg.Data, &taskCommand)

		if taskCommand.Command == "pause" || taskCommand.Command == "cancel" {
			s.cancelCh <- taskCommand
			s.store.UpdateStatus(taskCommand.ID, model.TaskStatus{
				Status: model.TaskState_Paused,
			})
		} else if taskCommand.Command == "resume" {
			// 获取任务，且状态为暂停
			// s.dispatch <- taskMsg
			task := s.store.GetTask(taskCommand.ID)
			s.dispatch <- task.Msg

		}

	})
}

func (s *worker) Stop() {
	s.syncer.shutdown()
	s.processor.Stop()
	s.subscriber.Shutdown()
}
