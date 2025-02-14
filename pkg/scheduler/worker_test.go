package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/google/uuid"
)

func DemoFunc(ctx TaskContext, t *model.Task) error {

	// load check point
	//
	var step int = 1
	checkPoint, ok := ctx.LoadCheckPoint().(int)
	if ok {
		step = checkPoint
	}

	for i := step; i <= 10; i++ {
		select {
		case <-ctx.Done():
			slog.Warn("task canceled")
			ctx.SaveCheckPoint(i)
			// save checkpoint
			return nil
		default:
			ctx.UpdateProgress(int32(i * 10))
			time.Sleep(time.Second / 2)
		}
	}
	// panic("crahs")

	return nil
}

func TestWorker(t *testing.T) {
	bus, err := eventbus.NewNatsBus(eventbus.NatsConfig{})
	if err != nil {
		t.Error(err)
		return
	}

	mux := NewTaskMux("worker")
	mux.HandleFunc("demo", func(ctx context.Context, t *model.Task) error {

		dctx, ok := ctx.(TaskContext)
		if !ok {
			return fmt.Errorf("not context")
		}

		//
		return DemoFunc(dctx, t)

	})

	worker := NewWorker(bus, mux)
	fmt.Println("work start")
	worker.Start()

	fmt.Println("work started")

	task := model.TaskMessage{}
	task.ID = uuid.NewString()
	task.Type = "demo"
	bus.Publish("tasks.queues", task)
	// task2 := TaskMessage{}
	// task2.ID = uuid.NewString()
	// task2.Type = "demo"
	// bus.Publish("tasks.queues", task2)
	time.Sleep(time.Second * 2)
	taskitem := worker.store.GetTask(task.ID)
	fmt.Println("---xxxxxxxxxx-----task status:", taskitem.Msg.ID, "----", taskitem.Status.Status, "----", taskitem.Progress.Progress)

	bus.Publish("tasks.commands", model.TaskCommand{
		ID:      task.ID,
		Command: "pause",
	})

	time.Sleep(time.Second * 5)
	bus.Publish("tasks.commands", model.TaskCommand{
		ID:      task.ID,
		Command: "resume",
	})

	time.Sleep(time.Second * 7)
	// worker.Start()
}
