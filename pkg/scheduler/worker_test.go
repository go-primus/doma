package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/google/uuid"
)

func TestWorker(t *testing.T) {
	bus, err := eventbus.NewNatsBus(eventbus.NatsConfig{})
	if err != nil {
		t.Error(err)
		return
	}
	worker := NewWorker(bus)
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
	time.Sleep(time.Second * 7)
	// worker.Start()
}
