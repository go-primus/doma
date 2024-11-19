package asynctask

import (
	"context"
	"fmt"

	"github.com/primus/primus/core/common/task"
)

type Task struct {
	typename string
	payload  []byte
	metadata *task.TaskMetadata

	// w is the ResultWriter for the task.
	w *ResultWriter
}

func (t *Task) Type() string                 { return t.typename }
func (t *Task) Payload() []byte              { return t.payload }
func (t *Task) Metadata() *task.TaskMetadata { return t.metadata }

// NewTask returns a new Task given a type name and payload data.
// Options can be passed to configure task processing behavior.
func NewTask(typename string, payload []byte) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
	}
}

func NewTaskWithMetadata(typename string, payload []byte, meta *task.TaskMetadata) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
		metadata: meta,
	}
}

// newTask creates a task with the given typename, payload and ResultWriter.
func newTask(typename string, payload []byte, w *ResultWriter) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
		w:        w,
	}
}

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
