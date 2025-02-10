package scheduler

var _ SchedulePolicy = (*fifoPolicy)(nil)

type fifoPolicy struct {
}

// Len implements SchedulePolicy.
func (f *fifoPolicy) Len() int {
	panic("unimplemented")
}

// Pop implements SchedulePolicy.
func (f *fifoPolicy) Pop() Task {
	panic("unimplemented")
}

// Push implements SchedulePolicy.
func (f *fifoPolicy) Push(task Task) (int, error) {
	panic("unimplemented")
}
