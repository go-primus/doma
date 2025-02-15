package core

import (
	"context"
	"fmt"

	"github.com/go-primus/doma/pkg/scheduler/internal/model"
)

// ////////////////////////
type TaskContext interface {
	context.Context
	UpdateProgress(progress int32)
	SaveCheckPoint(checkPoint any)
	LoadCheckPoint() any
	WriteResult(data any) error
}

type ProgressFunc func(progress int32)
type SaveCheckPointFunc func(checkPoint any)
type LoadCheckPointFunc func() any
type WriteResultFunc func(data any) error

type taskCtx struct {
	context.Context

	task model.TaskMessage

	updateProgressFunc ProgressFunc
	saveCheckPointFunc SaveCheckPointFunc
	loadCheckPointFunc LoadCheckPointFunc
	writeResultFunc    WriteResultFunc
}

func NewContext(ctx context.Context, msg model.TaskMessage, opts ...OptionFunc) TaskContext {
	tctx := &taskCtx{
		Context: ctx,
		task:    msg,
	}

	for _, opt := range opts {
		opt(tctx)
	}

	return tctx
}

// WriteResult implements TaskContext.
func (ctx *taskCtx) WriteResult(data any) error {
	if ctx.writeResultFunc == nil {
		return fmt.Errorf("unimplement")
	}
	return ctx.writeResultFunc(data)
}

// GetCheckPoint implements TaskContext.
func (ctx *taskCtx) LoadCheckPoint() any {
	if ctx.loadCheckPointFunc == nil {
		return fmt.Errorf("unimplement")
	}
	return ctx.loadCheckPointFunc()
}

// SaveCheckPoint implements TaskContext.
func (ctx *taskCtx) SaveCheckPoint(checkPoint any) {
	if ctx.saveCheckPointFunc == nil {
		return
	}
	ctx.saveCheckPointFunc(checkPoint)
}

func (ctx *taskCtx) UpdateProgress(progress int32) {
	if ctx.updateProgressFunc == nil {
		return
	}
	ctx.updateProgressFunc(progress)
}
