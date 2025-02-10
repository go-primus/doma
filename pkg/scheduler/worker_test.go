package scheduler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/google/uuid"
)

func DemoFunc(ctx context.Context, t *Task) error {

	updateProgress := func(progress int32) {
		taskCtx, ok := ctx.(TaskContext)
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
}

func TestWorker(t *testing.T) {
	bus, err := eventbus.NewNatsBus(eventbus.NatsConfig{})
	if err != nil {
		t.Error(err)
		return
	}

	mux := NewTaskMux("worker")
	mux.HandleFunc("demo", DemoFunc)

	worker := NewWorker(bus, mux)
	fmt.Println("work start")
	worker.Start()

	fmt.Println("work started")

	task := TaskMessage{}
	task.ID = uuid.NewString()
	task.Type = "demo"
	bus.Publish("tasks.queues", task)
	// task2 := TaskMessage{}
	// task2.ID = uuid.NewString()
	// task2.Type = "demo"
	// bus.Publish("tasks.queues", task2)
	time.Sleep(time.Second * 3)
	status := worker.store.GetStatus(task.ID)
	fmt.Println("---xxxxxxxxxx-----task status:", status.ID, "----", status.Status, "----", status.Progress)
	time.Sleep(time.Second * 7)
	// worker.Start()
}
