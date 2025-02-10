package scheduler

import (
	"context"
	"fmt"
	"log/slog"
)

type Processor interface {
	Start()
	Stop()
}

var _ Processor = (*processor)(nil)

type processor struct {
	queue <-chan TaskMessage

	handler Handler

	sema chan struct{}

	done chan struct{}
	quit chan struct{}

	//
	sync  chan<- *SyncMessage
	store TaskStore
}

func NewProcessor(queue <-chan TaskMessage) *processor {
	return &processor{
		// emitter: eventemitter.NewEventEmitter(),
		// bus:     bus,
		done:    make(chan struct{}),
		quit:    make(chan struct{}),
		handler: NotFoundHandler(),
		sema:    make(chan struct{}),
		queue:   queue,
	}
}

// Start implements Processor.
func (p *processor) Start() {
	go func() {
		fmt.Println("11111111")
		p.start()

		fmt.Println("22222222")
	}()
}

func (p *processor) start() {
	for {
		select {
		case <-p.done:
			return
		default:
			fmt.Println("exec")
			p.exec()
		}
	}
	// p.exec()
}

// Stop implements Processor.
func (p *processor) Stop() {

	close(p.quit)
	p.done <- struct{}{}
	//

}

func (p *processor) exec() {
	fmt.Println("execxxx")
	select {
	case <-p.quit:
		return
	// case p.sema <- struct{}{}: // acquire token
	case msg := <-p.queue:

		go func() {
			defer func() {
				// <-p.sema // release token
			}()
			p.execute(msg)
		}()
	}

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
