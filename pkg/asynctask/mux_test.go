package asynctask

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/reugn/go-quartz/quartz"
)

type fakeJob struct {
}

// Description implements quartz.Job.
func (f *fakeJob) Description() string {
	return fmt.Sprintf("fakejob")
}

// Execute implements quartz.Job.
func (f *fakeJob) Execute(ctx context.Context) error {
	fmt.Println("fake job execute")
	return nil
}

func makeFakeJob() quartz.Job {

	return &fakeJob{}
}

func makeFakeHandler() Handler {
	return HandlerFunc(
		func(ctx context.Context, task *Task) error {
			fmt.Println("----handle task----", task.Type())
			defer fmt.Println("---handle task -", task.Type(), "---end")
			time.Sleep(time.Second)
			return nil
		},
	)
}

func wrapHandle[T any](fn func(taskType string, payload T) error) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, task *Task) error {

		payload := new(T)
		err := json.Unmarshal(task.payload, payload)
		if err != nil {
			fmt.Println("unmarshal err:", err)
			return err
		}

		return fn(task.Type(), *payload)

	})
}

type Payload struct {
	Name string
	Size string
}

func TestMux(t *testing.T) {

	mux := NewTaskMux("taskmux")
	mux.Handle("xxxx", makeFakeHandler())

	mux.HandleFunc("yyyy", wrapHandle(func(typ string, payload Payload) error {
		fmt.Println("---", typ, "----", payload)
		return nil
	}))

	payload := Payload{}
	payload.Name = "namexxx"
	payload.Size = "sizexxx"
	bb, _ := json.Marshal(payload)

	mux.ProcessTask(context.Background(), &Task{

		typename: "yyyy",
		payload:  bb,
	})
}
