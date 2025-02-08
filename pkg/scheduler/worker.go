package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-primus/doma/pkg/eventbus"
)

type worker struct {
	processor processor
	syncer    syncer
}

func NewWorker(bus eventbus.EventBus) *worker {

	store := &taskStore{
		tasks: make(map[string]*TaskItem),
	}

	syncCh := make(chan *SyncMessage)
	processor := NewProcessor(bus)
	processor.sync = syncCh
	processor.store = store
	processor.handler = HandlerFunc(func(ctx context.Context, t *Task) error {

		updateProgress := func(progress int32) {
			taskCtx, ok := ctx.(taskCtx)
			if !ok {
				return
			}
			taskCtx.UpdateProgress(progress)
		}

		for i := 1; i <= 10; i++ {
			updateProgress(int32(i * 10))
			time.Sleep(time.Second / 2)
		}

		return nil
	})

	// processor.handler = fun

	syncer := syncer{}
	syncer.sync = syncCh
	syncer.store = store
	syncer.syncFunc = func(status TaskStatus) error {
		slog.Info("sync task", "task", status.ID, "status", status.Status, "---", status.Progress)
		return nil
	}

	return &worker{
		processor: *processor,
		syncer:    syncer,
	}
}

func (s *worker) Start() {
	s.syncer.Start()
	s.processor.Start()
}
