package retry

import "context"

// CheckCtxValid check if the context is valid
func CheckCtxValid(ctx context.Context) bool {
	return ctx.Err() != context.DeadlineExceeded && ctx.Err() != context.Canceled
}
