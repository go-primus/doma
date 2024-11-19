package asynctask

import (
	"context"
	"testing"
	"time"

	"github.com/jiyeyuran/go-eventemitter"
	"github.com/reugn/go-quartz/quartz"
)

func TestProcessor(t *testing.T) {
	processor := newProcessor()
	processor.emitter = eventemitter.NewEventEmitter()
	processor.handler = makeFakeHandler()
	processor.onnn()

	processor.enqueue("task1")
	processor.enqueue("task2")
	processor.enqueue("task3")
	processor.enqueue("task4")

	time.Sleep(time.Second * 3)
}

type fakeTrigger struct {
}

// func makeFakeTrigger() quartz.tr

func TestQuartz(t *testing.T) {
	sched := quartz.NewStdScheduler()
	sched.Start(context.Background())
	sched.ScheduleJob(quartz.NewJobDetail(makeFakeJob(), quartz.NewJobKey("")), quartz.NewRunOnceTrigger(0))

	time.Sleep(time.Second * 3)
}
