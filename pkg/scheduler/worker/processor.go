package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-primus/doma/pkg/scheduler/core"
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

type processorParams struct {
	dispatch     <-chan model.TaskMessage
	cancelations *cancelation.Cancelations
	sync         chan<- *model.SyncMessage
	handler      Handler
	store        store.TaskStore
}

func newProcessor(params processorParams) *processor {
	p := &processor{
		dispatch:     params.dispatch,
		done:         make(chan struct{}),
		quit:         make(chan struct{}),
		sema:         make(chan struct{}),
		handler:      NotFoundHandler(),
		cancelations: params.cancelations,
		//
		sync:  params.sync,
		store: params.store,
	}

	if params.handler != nil {
		p.handler = params.handler
	}
	return p
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

	// check context before starting a worker goroutine.
	select {
	case <-ctx.Done():
		// alreay canceled (e.g. deadline exceeded).
		p.handleFailedMessage(msg, ctx.Err())
		return
	default:
	}

	resCh := make(chan error, 1)
	go func() {

		taskCtx := core.NewContext(ctx, msg,
			core.WithProgressFunc(func(progress int32) {
				p.handleProgressMessage(msg, progress)
			}),
			core.WithSaveCheckPointFunc(func(checkPoint any) {
				p.saveCheckPoint(msg.ID, checkPoint)
			}),
			core.WithLoadCheckPointFunc(func() any {
				return p.loadCheckPoint(msg.ID)
			}),
		)

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

///////////////////////
// processor handler //
///////////////////////

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
