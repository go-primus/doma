package scheduler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/eventbus"
)

type worker struct {
	bus eventbus.EventBus

	dispatch chan<- TaskMessage // dispatch queue

	//
	processor processor
	syncer    syncer

	//
	finished chan<- TaskMessage

	//
	store TaskStore
	// local scheduler

}

func NewWorker(bus eventbus.EventBus, handler Handler) *worker {

	store := &taskStore{
		tasks: make(map[string]*TaskItem),
	}

	dispatch := make(chan TaskMessage, 10)
	syncCh := make(chan *SyncMessage)

	// mux := NewTaskMux("worker")

	processor := NewProcessor(dispatch)
	processor.sync = syncCh
	processor.handler = handler

	syncer := newSyncer()
	syncer.sync = syncCh
	syncer.store = store

	syncer.syncFunc = func(syncmsg SyncMessage) error {
		slog.Info("sync task", "task", syncmsg.task.ID, "status", syncmsg.status.Status, "---", syncmsg.progress.Progress)
		return nil
	}

	return &worker{
		bus:       bus,
		dispatch:  dispatch,
		processor: *processor,
		store:     store,
		syncer:    *syncer,
	}
}

func (s *worker) Start() {
	s.syncer.Start()
	s.processor.Start()
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

	})
}

func (s *worker) Stop() {
	s.syncer.shutdown()
	s.processor.Stop()
}
