package scheduler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/go-primus/doma/pkg/scheduler/internal/cancelation"
)

type worker struct {
	bus eventbus.EventBus

	dispatch chan<- TaskMessage // dispatch queue

	//
	processor  *processor
	syncer     *syncer
	subscriber *subscriber

	//
	cancelCh chan<- TaskCommand

	//
	finished chan<- TaskMessage

	//
	store TaskStore
}

func NewWorker(bus eventbus.EventBus, handler Handler) *worker {

	store := &taskStore{
		tasks: make(map[string]*TaskItem),
	}

	dispatch := make(chan TaskMessage, 10)
	syncCh := make(chan *SyncMessage)

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

	syncer.syncFunc = func(syncmsg SyncMessage) error {
		slog.Info("sync task", "task", syncmsg.task.ID, "status", syncmsg.status.Status, "---", syncmsg.progress.Progress)
		return nil
	}

	cancelCh := make(chan TaskCommand)

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
		var taskMsg TaskMessage
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
		var taskCommand TaskCommand
		json.Unmarshal(msg.Data, &taskCommand)

		if taskCommand.Command == "pause" || taskCommand.Command == "cancel" {
			s.cancelCh <- taskCommand
		} else if taskCommand.Command == "resume" {
			// 获取任务，且状态为暂停
			// s.dispatch <- taskMsg
			task := s.store.GetTask(taskCommand.ID)
			s.dispatch <- task.msg

		}

	})
}

func (s *worker) Stop() {
	s.syncer.shutdown()
	s.processor.Stop()
}
