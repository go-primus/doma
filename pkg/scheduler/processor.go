package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/scheduler/internal/cancelation"
	"github.com/go-primus/doma/pkg/scheduler/internal/errors"
	"github.com/go-primus/doma/pkg/scheduler/internal/model"
	"github.com/go-primus/doma/pkg/scheduler/internal/store"
)

type Processor interface {
	Start()
	Stop()
}

var _ Processor = (*processor)(nil)

type processor struct {
	dispatch <-chan model.TaskMessage

	handler Handler

	store store.TaskStore

	sema chan struct{}

	done chan struct{}
	quit chan struct{}

	cancelations *cancelation.Cancelations

	//
	sync chan<- *model.SyncMessage // report sync message ,task status
	// finished chan<- TaskMessage
}

func NewProcessor(dispatch <-chan model.TaskMessage, cancelations *cancelation.Cancelations) *processor {
	return &processor{
		dispatch:     dispatch,
		done:         make(chan struct{}),
		quit:         make(chan struct{}),
		sema:         make(chan struct{}),
		handler:      NotFoundHandler(),
		cancelations: cancelations,
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

func (p *processor) execute(msg model.TaskMessage) {

	ctx, cancel := context.WithCancel(context.Background())
	p.cancelations.Add(msg.ID, cancel)
	defer func() {
		cancel()
		p.cancelations.Delete(msg.ID)
	}()

	taskCtx := NewContext(ctx, msg, WithProgressFunc(func(progress int32) {
		p.handleProgressMessage(msg, progress)
	}), WithSaveCheckPointFunc(func(checkPoint any) {
		p.saveCheckPoint(msg.ID, checkPoint)
	}), WithLoadCheckPointFunc(func() any {
		return p.loadCheckPoint(msg.ID)
	}))

	select {
	case <-ctx.Done():
		return
	default:
	}

	resCh := make(chan error, 1)
	go func() {
		task := &model.Task{
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

func (p *processor) perform(ctx context.Context, task *model.Task) (err error) {

	defer func() {
		if x := recover(); x != nil {
			errMsg := fmt.Sprintf("panic: %v", x)
			slog.Error(errMsg)
			err = &errors.PanicError{
				ErrMsg: errMsg,
			}
		}
	}()
	return p.handler.ProcessTask(ctx, task)
}

func (p *processor) handleFailedMessage(msg model.TaskMessage, err error) {

	taskStatus := model.TaskStatus{}

	if errors.Is(err, context.Canceled) {
		return
	}
	taskStatus.Status = model.TaskState_Failed
	taskStatus.Err = err

	p.sync <- &model.SyncMessage{
		Task:   msg,
		Status: taskStatus,
	}
}

func (p *processor) handleSucceededMessage(msg model.TaskMessage) {

	taskStatus := model.TaskStatus{}
	taskStatus.Status = model.TaskState_Completed

	p.sync <- &model.SyncMessage{
		Task:   msg,
		Status: taskStatus,
	}

}

func (p *processor) handleProgressMessage(msg model.TaskMessage, progress int32) {

	taskStatus := model.TaskStatus{}
	taskStatus.Status = model.TaskState_Running
	taskProgress := model.TaskProgress{
		Progress: int(progress),
	}

	p.sync <- &model.SyncMessage{
		Task:     msg,
		Status:   taskStatus,
		Progress: taskProgress,
	}
}

func (p *processor) saveCheckPoint(taskId string, checkpoint any) {
	p.store.SaveCheckPoint(taskId, checkpoint)
}

func (p *processor) loadCheckPoint(taskId string) any {

	return p.store.LoadCheckPoint(taskId)
}
