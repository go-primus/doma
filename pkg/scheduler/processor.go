package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/eventbus"
)

type Processor interface {
	Start()
	Stop()
}

var _ Processor = (*processor)(nil)

type processor struct {
	// emitter eventemitter.IEventEmitter
	bus eventbus.EventBus

	handler Handler

	done chan struct{}
	quit chan struct{}

	//
	sync  chan<- *SyncMessage
	store TaskStore
}

func NewProcessor(bus eventbus.EventBus) *processor {
	return &processor{
		// emitter: eventemitter.NewEventEmitter(),
		bus:     bus,
		done:    make(chan struct{}),
		quit:    make(chan struct{}),
		handler: NotFoundHandler(),
	}
}

// Start implements Processor.
func (p *processor) Start() {
	// go func() {
	p.start()
	// }()
}

func (p *processor) start() {
	// for {
	// 	select {
	// 	case <-p.done:
	// 		return
	// 	default:
	// 		p.exec()
	// 	}
	// }
	p.exec()
}

// Stop implements Processor.
func (p *processor) Stop() {
	//
	// p.emitter.RemoveAllListeners("")

	close(p.quit)
	p.done <- struct{}{}
	//

}

func (p *processor) exec() {
	// select {
	// case <-p.quit:
	// 	return
	// case p.sema <- struct{}{}: // acquire token
	// }

	p.bus.Subscribe("tasks.queues", func(msg *eventbus.Msg) {
		var taskMsg TaskMessage
		json.Unmarshal(msg.Data, &taskMsg)

		fmt.Println("handle task queues msg")
		//
		p.execute(taskMsg)

	})

	// subscribe tasks.command

}

func (p *processor) execute(msg TaskMessage) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	taskCtx := taskCtx{
		Context: ctx,
		task:    msg,
		updateProgressFunc: func(progress int32) {
			p.handleProgressMessage(msg, progress)
		},
	}

	select {
	case <-ctx.Done():
		return
	default:
	}

	// check task
	// 查询任务是否正在执行，如果正在执行忽略
	//
	p.store.AddTask(msg)

	// update task status

	resCh := make(chan error, 1)
	go func() {
		task := &Task{
			ID:      msg.ID,
			Type:    msg.Type,
			Payload: msg.Payload,
		}
		resCh <- p.perform(taskCtx, task)
	}()

	select {
	case <-ctx.Done():
		p.handleFailedMessage(msg, ctx.Err())
		return
	case resErr := <-resCh:
		if resErr != nil {
			p.handleFailedMessage(msg, resErr)
			return
		}
		p.handleSucceededMessage(msg)
	}
}

func (p *processor) perform(ctx context.Context, task *Task) error {

	defer func() {
		if x := recover(); x != nil {
			errMsg := fmt.Sprintf("panic: %v", x)
			slog.Error(errMsg)
		}
	}()
	return p.handler.ProcessTask(ctx, task)
}

func (p *processor) handleFailedMessage(msg TaskMessage, err error) {
	taskStatus := TaskStatus{}
	taskStatus.ID = msg.ID
	taskStatus.Status = "failed"
	taskStatus.Err = err.Error()

	p.store.UpdateStatus(msg.ID, taskStatus)

	// p.sync <- &SyncMessage{
	// 	task:   msg,
	// 	status: taskStatus,
	// }
	// p.bus.Publish("tasks.updates", taskStatus)
}

func (p *processor) handleSucceededMessage(msg TaskMessage) {

	taskStatus := TaskStatus{}
	taskStatus.ID = msg.ID
	taskStatus.Status = "completed"
	taskStatus.Progress = 100

	p.store.UpdateStatus(msg.ID, taskStatus)

	// p.sync <- &SyncMessage{
	// 	task:   msg,
	// 	status: taskStatus,
	// }

	// 失败需要重发
	// p.bus.Publish("tasks.updates", taskStatus)
}

func (p *processor) handleProgressMessage(msg TaskMessage, progress int32) {

	taskStatus := TaskStatus{}
	taskStatus.ID = msg.ID
	taskStatus.Status = "running"
	taskStatus.Progress = int(progress)

	p.store.UpdateStatus(msg.ID, taskStatus)

	// p.sync <- &SyncMessage{
	// 	task:   msg,
	// 	status: taskStatus,
	// }
}
