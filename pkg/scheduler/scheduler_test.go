package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestScheduler(t *testing.T) {

	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {

		t.Error(err)
		return
	}
	sched := scheduler{
		nc: conn,
		store: &taskStore{
			tasks: make(map[string]TaskStatus),
		},
	}
	sched.Start()
	id, err := sched.SubmitTask(Task{
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
