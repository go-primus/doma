package task

import (
	"context"
	"encoding/json"
	"fmt"
)

type TaskContext interface {
	context.Context
	TaskID() string
	TaskType() string
	TaskProgress
}

type TaskProgress interface {
	PublishProgress(success, total int32) error
}

type emptyTaskProgress struct {
}

func NewEmptyTaskProgress() TaskProgress {
	return &emptyTaskProgress{}
}

func (p *emptyTaskProgress) PublishProgress(success, total int32) error {

	return nil
}

type taskContext struct {
	context.Context
	taskId   string
	taskType string
	// w is the ResultWriter for the task.
	w *ResultWriter
}

func (t *taskContext) TaskID() string {
	return t.taskId
}

func (t *taskContext) TaskType() string {
	return t.taskType
}

// PublishProgress implements TaskContext.
func (t *taskContext) PublishProgress(success int32, total int32) error {

	if t.w != nil {
		bb, _ := json.Marshal(struct {
			Success int32
			Total   int32
		}{
			Success: success,
			Total:   total,
		})
		t.w.Write(bb)
	}
	return nil

}

func WrapTaskContext(ctx context.Context, taskId string, taskType string, w *ResultWriter) TaskContext {
	return &taskContext{
		Context:  ctx,
		taskId:   taskId,
		taskType: taskType,
		w:        w,
	}
}

// //

// ResultWriter is a client interface to write result data for a task.
// It writes the data to the redis instance the server is connected to.
type ResultWriter struct {
	id    string // task ID this writer is responsible for
	qname string // queue name the task belongs to
	// broker base.Broker
	ctx context.Context // context associated with the task
}

// Write writes the given data as a result of the task the ResultWriter is associated with.
func (w *ResultWriter) Write(data []byte) (n int, err error) {
	select {
	case <-w.ctx.Done():
		return 0, fmt.Errorf("failed to result task result: %v", w.ctx.Err())
	default:
	}
	// return w.broker.WriteResult(w.qname, w.id, data)
	return 0, nil
}

// TaskID returns the ID of the task the ResultWriter is associated with.
func (w *ResultWriter) TaskID() string {
	return w.id
}
