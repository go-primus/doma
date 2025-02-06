package svc

import (
	"context"
	"encoding/json"

	"github.com/go-primus/doma/core/common/task"
	"github.com/go-primus/doma/pkg/asynctask"
)

type TaskMux struct {
	s     Service
	_name string
	asynctask.TaskMux
}

func NewTaskMux(muxName string) *TaskMux {
	mux := new(TaskMux)
	mux._name = muxName
	mux.TaskMux = *asynctask.NewTaskMux(muxName)
	return mux
}

func NewServiceTaskMux(s Service) *TaskMux {
	mux := NewTaskMux(s.Name())
	mux.s = s
	return mux
}

func (mux *TaskMux) Subscribe() {
	s := mux.s

	if s == nil {
		return
	}
	s.SubscribeTask(mux.OnTask)
}

func (mux *TaskMux) OnTask(ctx context.Context, task task.Task) (any, error) {

	bb, _ := json.Marshal(task.Payload)
	asyncTask := asynctask.NewTaskWithMetadata(task.TaskType, bb, task.Metadata)
	// asyncTask.
	err := mux.ProcessTask(ctx, asyncTask)

	return nil, err
}
