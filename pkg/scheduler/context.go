package scheduler

import (
	"context"
)

// ////////////////////////
type TaskContext interface {
	context.Context
	UpdateProgress(progress int32)
	SaveCheckPoint(checkPoint any)
	LoadCheckPoint() any
}

type ProgressFunc func(progress int32)
type SaveCheckPointFunc func(checkPoint any)
type LoadCheckPointFunc func() any

type taskCtx struct {
	context.Context

	task TaskMessage

	updateProgressFunc ProgressFunc
	saveCheckPointFunc SaveCheckPointFunc
	loadCheckPointFunc LoadCheckPointFunc
}

func NewContext(ctx context.Context, msg TaskMessage, opts ...OptionFunc) TaskContext {
	tctx := &taskCtx{
		Context: ctx,
		task:    msg,
	}

	for _, opt := range opts {
		opt(tctx)
	}

	return tctx
}

// GetCheckPoint implements TaskContext.
func (ctx *taskCtx) LoadCheckPoint() any {
	return ctx.loadCheckPointFunc()
}

// SaveCheckPoint implements TaskContext.
func (ctx *taskCtx) SaveCheckPoint(checkPoint any) {
	ctx.saveCheckPointFunc(checkPoint)
}

func (ctx *taskCtx) UpdateProgress(progress int32) {
	ctx.updateProgressFunc(progress)
}
