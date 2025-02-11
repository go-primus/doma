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
	processor.store = store
	processor.handler = handler

	syncer := syncer{}
	syncer.sync = syncCh
	syncer.store = store
	syncer.syncFunc = func(status TaskStatus) error {
		slog.Info("sync task", "task", status.ID, "status", status.Status, "---", status.Progress)
		return nil
	}

	return &worker{
		bus:       bus,
		dispatch:  dispatch,
		processor: *processor,
		store:     store,
		syncer:    syncer,
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
		//
		s.dispatch <- taskMsg
		fmt.Println("handle task queues msgxxx")
	})
}

func (s *worker) Stop() {
	s.syncer.shutdown()
	s.processor.Stop()
}
