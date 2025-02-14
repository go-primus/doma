package scheduler

import "github.com/go-primus/doma/pkg/scheduler/internal/model"

var _ SchedulePolicy = (*fifoPolicy)(nil)

type fifoPolicy struct {
}

// Len implements SchedulePolicy.
func (f *fifoPolicy) Len() int {
	panic("unimplemented")
}

// Pop implements SchedulePolicy.
func (f *fifoPolicy) Pop() model.Task {
	panic("unimplemented")
}

// Push implements SchedulePolicy.
func (f *fifoPolicy) Push(task model.Task) (int, error) {
	panic("unimplemented")
}
