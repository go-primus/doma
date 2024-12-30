package retry

import (
	"context"
	"log/slog"
	"time"

	"github.com/cockroachdb/errors"
)

// Do will run function with retry mechanism.
// fn is the func to run.
// Option can control the retry times and timeout.
func Do(ctx context.Context, fn func() error, opts ...Option) error {
	if !CheckCtxValid(ctx) {
		return ctx.Err()
	}

	log := slog.Default()
	c := newDefaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	var lastErr error

	for i := uint(0); c.attempts == 0 || i < c.attempts; i++ {
		if err := fn(); err != nil {
			if i%4 == 0 {
				log.Warn("retry func failed", "retried", i, "err", err)
			}

			if !IsRecoverable(err) {
				isContextErr := errors.IsAny(err, context.Canceled, context.DeadlineExceeded)
				log.Warn("retry func failed, not be recoverable",
					"retried", i,
					"attempt", c.attempts,
					"isContextErr", isContextErr,
				)
				if isContextErr && lastErr != nil {
					return lastErr
				}
				return err
			}
			if c.isRetryErr != nil && !c.isRetryErr(err) {
				log.Warn("retry func failed, not be retryable",
					"retried", i,
					"attempt", c.attempts,
				)
				return err
			}

			deadline, ok := ctx.Deadline()
			if ok && time.Until(deadline) < c.sleep {
				isContextErr := errors.IsAny(err, context.Canceled, context.DeadlineExceeded)
				log.Warn("retry func failed, deadline",
					"retried", i,
					"attempt", c.attempts,
					"isContextErr", isContextErr,
				)
				if isContextErr && lastErr != nil {
					return lastErr
				}
				return err
			}

			lastErr = err

			select {
			case <-time.After(c.sleep):
			case <-ctx.Done():
				log.Warn("retry func failed, ctx done",
					"retried", i,
					"attempt", c.attempts,
				)
				return lastErr
			}

			c.sleep *= 2
			if c.sleep > c.maxSleepTime {
				c.sleep = c.maxSleepTime
			}
		} else {
			return nil
		}
	}
	if lastErr != nil {
		log.Warn("retry func failed, reach max retry",
			"attempt", c.attempts,
		)
	}
	return lastErr
}

// Do will run function with retry mechanism.
// fn is the func to run, return err and shouldRetry flag.
// Option can control the retry times and timeout.
func Handle(ctx context.Context, fn func() (bool, error), opts ...Option) error {
	if !CheckCtxValid(ctx) {
		return ctx.Err()
	}

	log := slog.Default()
	c := newDefaultConfig()

	for _, opt := range opts {
		opt(c)
	}

	var lastErr error
	for i := uint(0); i < c.attempts; i++ {
		if shouldRetry, err := fn(); err != nil {
			if i%4 == 0 {
				log.Warn("retry func failed", "retried", i, "err", err)
			}

			if !shouldRetry {
				if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) && lastErr != nil {
					return lastErr
				}
				return err
			}

			deadline, ok := ctx.Deadline()
			if ok && time.Until(deadline) < c.sleep {
				// to avoid sleep until ctx done
				if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) && lastErr != nil {
					return lastErr
				}
				return err
			}

			lastErr = err

			select {
			case <-time.After(c.sleep):
			case <-ctx.Done():
				return lastErr
			}

			c.sleep *= 2
			if c.sleep > c.maxSleepTime {
				c.sleep = c.maxSleepTime
			}
		} else {
			return nil
		}
	}
	return lastErr
}

// errUnrecoverable is error instance for unrecoverable.
var errUnrecoverable = errors.New("unrecoverable error")

// Unrecoverable method wrap an error to unrecoverableError. This will make retry
// quick return.
func Unrecoverable(err error) error {
	return err
}

// IsRecoverable is used to judge whether the error is wrapped by unrecoverableError.
func IsRecoverable(err error) bool {
	return !errors.Is(err, errUnrecoverable)
}
