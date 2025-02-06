package job

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

var _ Job = (*DemoJob)(nil)

type DemoJob struct {
	BaseJob
	fn func() error
}

func NewFunc(fn func() error) *DemoJob {
	return &DemoJob{
		BaseJob: *NewBaseJob(context.Background(), 11, 11),
		fn:      fn,
	}
}

// Execute implements Job.
func (d *DemoJob) Execute() error {

	slog.Info("execute demo job")
	return d.fn()
}

func (d *DemoJob) PostExecute() error {
	return nil
}

func TestJob(t *testing.T) {

	js := NewScheduler()
	js.Start()

	// js.Add(&DemoJob{
	// 	BaseJob: *NewBaseJob(context.Background(), 11, 11),
	// })

	fn := NewFunc(func() error {

		time.Sleep(time.Second)
		slog.Info("func run")
		return nil
	})
	js.Add(fn)

	fn.Wait()
	slog.Info("----")

	time.Sleep(time.Second * 10)

}
