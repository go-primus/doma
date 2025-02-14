package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/go-primus/doma/pkg/scheduler/internal/store"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var _ Scheduler = (*scheduler)(nil)

type scheduler struct {
	nc *nats.Conn

	queue SchedulePolicy
	store store.TaskStore
}

// Add implements Scheduler.
func (s *scheduler) Add(task model.Task) error {
	panic("unimplemented")
}

// Stop implements Scheduler.
func (s *scheduler) Stop() error {
	panic("unimplemented")
}

// RetryTask implements Scheduler.
func (s *scheduler) RetryTask(id string) error {
	panic("unimplemented")
}

// DeleteTask implements Scheduler.
func (s *scheduler) DeleteTask(id string) error {
	panic("unimplemented")
}

// PauseTask implements Scheduler.
func (s *scheduler) PauseTask(id string) error {
	panic("unimplemented")
}

// ResumeTask implements Scheduler.
func (s *scheduler) ResumeTask(id string) error {
	panic("unimplemented")
}

// GetTask implements Scheduler.
func (s *scheduler) GetTask(id string) model.TaskStatus {
	status := s.store.GetStatus(id)
	return status
}

func (s *scheduler) Start() error {
	s.nc.Subscribe("tasks.updates", func(msg *nats.Msg) {
		var status model.TaskStatus
		json.Unmarshal(msg.Data, &status)
	})

	s.nc.Subscribe("tasks.command", func(msg *nats.Msg) {
		//
		var cmd model.TaskCommand
		json.Unmarshal(msg.Data, &cmd)

		switch cmd.Command {
		case "pause":
			s.store.PauseTask(cmd.ID)
		case "resume":
			s.store.ResumeTask(cmd.ID)
		}
		s.nc.Publish("tasks.worker.command", msg.Data)

	})
	s.work()
	return nil
}

func (s *scheduler) schedule() {
	// schedule  loop

	for {

	}
}

func (s *scheduler) SubmitTask(task model.Task) (string, error) {
	task.ID = uuid.NewString()
	// initialStatus := TaskStatus{
	// 	Status: "pending",
	// }
	// s.store.AddTask(task.ID, initialStatus)

	data, _ := json.Marshal(task)

	return task.ID, s.nc.Publish("tasks.queue", data)

}

func (s *scheduler) work() {

	s.nc.Subscribe("tasks.queue", func(msg *nats.Msg) {
		var task model.Task
		json.Unmarshal(msg.Data, &task)
		//

		//
		fmt.Println("handle task")

		// taskCtx := taskCtx{taskID}
		//err :=  processTask(taskCtx, task)
		// mux.processTask()

	})
}

func (s *scheduler) updateStatus(status model.TaskStatus) {
	data, _ := json.Marshal(status)
	s.nc.Publish("tasks.updates", data)
}
