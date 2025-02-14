package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/go-primus/doma/pkg/scheduler/internal/store"
	"github.com/nats-io/nats.go"
)

func TestScheduler(t *testing.T) {

	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {

		t.Error(err)
		return
	}
	sched := scheduler{
		nc:    conn,
		store: store.NewTaskStore(),
	}

	// worker := NewProcessor()
	// worker.handler = NewTaskMux("mux")
	// worker.Start()

	sched.Start()
	id, err := sched.SubmitTask(model.Task{
		Type:    "move",
		Payload: "",
	})
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	status := sched.GetTask(id)

	fmt.Println(status)
}
