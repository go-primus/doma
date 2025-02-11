package scheduler

import "context"

// ////////////////////////
type TaskContext interface {
	context.Context
	UpdateProgress(progress int32)
}

type ProgressFunc func(progress int32)

func NewContext(ctx context.Context, msg TaskMessage, progressFunc ProgressFunc) TaskContext {
	return &taskCtx{
		Context:            ctx,
		task:               msg,
		updateProgressFunc: progressFunc,
	}
}

type taskCtx struct {
	context.Context

	task               TaskMessage
	updateProgressFunc func(progress int32)
}

func (ctx *taskCtx) UpdateProgress(progress int32) {
	ctx.updateProgressFunc(progress)
}
