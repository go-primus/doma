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
	dispatch <-chan TaskMessage

	handler Handler

	sema chan struct{}

	done chan struct{}
	quit chan struct{}

	//
	sync chan<- *SyncMessage
	// finished chan<- TaskMessage
}

func NewProcessor(dispatch <-chan TaskMessage) *processor {
	return &processor{
		// emitter: eventemitter.NewEventEmitter(),
		// bus:     bus,
		done:     make(chan struct{}),
		quit:     make(chan struct{}),
		handler:  NotFoundHandler(),
		sema:     make(chan struct{}),
		dispatch: dispatch,
	}
}

// Start implements Processor.
func (p *processor) Start() {
	go func() {
		p.start()
	}()
}

func (p *processor) start() {
	for {
		select {
		case <-p.done:
			return
		default:
			p.exec()
		}
	}
}

// Stop implements Processor.
func (p *processor) Stop() {

	close(p.quit)
	p.done <- struct{}{}
	//

}

func (p *processor) exec() {
	select {
	case <-p.quit:
		return
	// case p.sema <- struct{}{}: // acquire token
	case msg := <-p.dispatch:

		go func() {
			defer func() {
				// p.finished <- msg
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

	taskCtx := NewContext(ctx, msg, func(progress int32) {
		p.handleProgressMessage(msg, progress)
	})

	select {
	case <-ctx.Done():
		return
	default:
	}

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
	taskStatus.Status = TaskState_Failed
	taskStatus.Err = err

	p.sync <- &SyncMessage{
		task:   msg,
		status: taskStatus,
	}
}

func (p *processor) handleSucceededMessage(msg TaskMessage) {

	taskStatus := TaskStatus{}
	taskStatus.Status = TaskState_Completed

	p.sync <- &SyncMessage{
		task:   msg,
		status: taskStatus,
	}

}

func (p *processor) handleProgressMessage(msg TaskMessage, progress int32) {

	taskStatus := TaskStatus{}
	taskStatus.Status = TaskState_Running
	taskProgress := TaskProgress{
		Progress: int(progress),
	}

	p.sync <- &SyncMessage{
		task:     msg,
		status:   taskStatus,
		progress: taskProgress,
	}
}
