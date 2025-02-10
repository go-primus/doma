package scheduler

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/eventbus"
)

type worker struct {
	bus       eventbus.EventBus
	queue     chan<- TaskMessage
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

	queue := make(chan TaskMessage, 10)
	syncCh := make(chan *SyncMessage)

	// mux := NewTaskMux("worker")

	processor := NewProcessor(queue)
	processor.sync = syncCh
	processor.store = store
	processor.handler = handler
	// processor.handler = HandlerFunc(func(ctx context.Context, t *Task) error {

	// 	updateProgress := func(progress int32) {
	// 		taskCtx, ok := ctx.(taskCtx)
	// 		if !ok {
	// 			return
	// 		}
	// 		taskCtx.UpdateProgress(progress)
	// 	}

	// 	for i := 1; i <= 10; i++ {
	// 		updateProgress(int32(i * 10))
	// 		time.Sleep(time.Second / 2)
	// 	}

	// 	return nil
	// })

	// processor.handler = fun

	syncer := syncer{}
	syncer.sync = syncCh
	syncer.store = store
	syncer.syncFunc = func(status TaskStatus) error {
		slog.Info("sync task", "task", status.ID, "status", status.Status, "---", status.Progress)
		return nil
	}

	return &worker{
		bus:       bus,
		queue:     queue,
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
		s.queue <- taskMsg
		fmt.Println("handle task queues msgxxx")
	})
}

func (s *worker) Stop() {
	s.syncer.shutdown()
	s.processor.Stop()
}
