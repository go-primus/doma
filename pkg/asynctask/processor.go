package asynctask

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/jiyeyuran/go-eventemitter"
	"github.com/primus/primus/pkg/asynctask/internal/base"
	"github.com/primus/primus/pkg/asynctask/internal/errors"
	"github.com/primus/primus/pkg/logger"
)

type processor struct {
	logger  logger.Logger
	broker  base.Broker
	emitter eventemitter.IEventEmitter

	handler Handler

	taskCheckInterval time.Duration

	shutdownTimeout time.Duration

	// channel via which to send sync requests to syncer.
	// syncRequestCh chan<- *syncRequest

	// rate limiter to prevent spamming logs with a bunch of errors.
	errLogLimiter *rate.Limiter

	// sema is a counting semaphore to ensure the number of active workers
	// does not exceed the limit.
	sema chan struct{}

	// channel to communicate back to the long running "processor" goroutine.
	// once is used to send value to the channel only once.
	done chan struct{}
	once sync.Once

	// quit channel is closed when the shutdown of the "processor" goroutine starts.
	quit chan struct{}

	// abort channel communicates to the in-flight worker goroutines to stop.
	abort chan struct{}
}

func newProcessor() *processor {
	return &processor{
		sema: make(chan struct{}),
	}
}

// Note: stops only the "processor" goroutine, does not stop workers.
// It's safe to call this method multiple times.
func (p *processor) stop() {
	p.once.Do(func() {
		p.logger.Debug("Processor shutting down...")
		// Unblock if processor is waiting for sema token.
		close(p.quit)
		// Signal the processor goroutine to stop processing tasks
		// from the queue.
		p.done <- struct{}{}
	})
}

// NOTE: once shutdown, processor cannot be re-started.
func (p *processor) shutdown() {
	p.stop()

	time.AfterFunc(p.shutdownTimeout, func() { close(p.abort) })

	p.logger.Info("Waiting for all workers to finish...")
	// block until all workers have released the token
	for i := 0; i < cap(p.sema); i++ {
		p.sema <- struct{}{}
	}
	p.logger.Info("All workers have finished")
}

func (p *processor) start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-p.done:
				p.logger.Debug("Processor done")
				return
			default:
				p.exec()
			}
		}
	}()
}

func (p *processor) enqueue(name string) {
	p.emitter.Emit("task", &base.TaskMessage{
		Type: name,
	})
}

func (p *processor) onnn() {
	p.emitter.On("task", func(msg *base.TaskMessage) {
		p.executeTask(msg)
	})
}

// exec pulls a task out of the queue and starts a worker goroutine to
// process the task.
func (p *processor) exec() {
	select {
	case <-p.quit:
		return
	case p.sema <- struct{}{}: // acquire token
		// qnames := p.queues()
		msg, _, err := p.broker.Dequeue()
		switch {
		case errors.Is(err, ErrNoProcessableTask):
			p.logger.Debug("All queues are empty")
			// Queues are empty, this is a normal behavior.
			// Sleep to avoid slamming redis and let scheduler move tasks into queues.
			// Note: We are not using blocking pop operation and polling queues instead.
			// This adds significant load to redis.
			time.Sleep(p.taskCheckInterval)
			<-p.sema // release token
			return
		case err != nil:
			if p.errLogLimiter.Allow() {
				p.logger.Errorf("Dequeue error: %v", err)
			}
			<-p.sema // release token
			return
		}

		p.executeTask(msg)

	}
}

func (p *processor) executeTask(msg *base.TaskMessage) {

	// lease := base.NewLease(leaseExpirationTime)
	// deadline := p.computeDeadline(msg)
	// p.starting <- &workerInfo{msg, time.Now(), deadline, lease}
	go func() {
		// defer func() {
		// 	p.finished <- msg
		// 	<-p.sema // release token
		// }()

		// ctx, cancel := asynqcontext.New(p.baseCtxFn(), msg, deadline)
		// p.cancelations.Add(msg.ID, cancel)
		// defer func() {
		// 	cancel()
		// 	p.cancelations.Delete(msg.ID)
		// }()
		ctx := context.Background()

		// check context before starting a worker goroutine.
		// select {
		// case <-ctx.Done():
		// 	// already canceled (e.g. deadline exceeded).
		// 	p.handleFailedMessage(ctx, lease, msg, ctx.Err())
		// 	return
		// default:
		// }

		resCh := make(chan error, 1)
		go func() {
			task := newTask(
				msg.Type,
				msg.Payload,
				nil,
				// &ResultWriter{
				// 	id:     msg.ID,
				// 	qname:  msg.Queue,
				// 	broker: p.broker,
				// 	ctx:    ctx,
				// },
			)
			resCh <- p.perform(ctx, task)
		}()

		select {
		// case <-p.abort:
		// 	// time is up, push the message back to queue and quit this worker goroutine.
		// 	p.logger.Warnf("Quitting worker. task id=%s", msg.ID)
		// 	p.requeue(lease, msg)
		// 	return
		// case <-lease.Done():
		// 	cancel()
		// 	p.handleFailedMessage(ctx, lease, msg, ErrLeaseExpired)
		// 	return
		case <-ctx.Done():
			p.handleFailedMessage(ctx, msg, ctx.Err())
			return
		case resErr := <-resCh:
			if resErr != nil {
				p.handleFailedMessage(ctx, msg, resErr)
				return
			}
			p.handleSucceededMessage(msg)
		}
	}()
}

// computeDeadline returns the given task's deadline,
// func (p *processor) computeDeadline(msg *base.TaskMessage) time.Time {
// 	if msg.Timeout == 0 && msg.Deadline == 0 {
// 		p.logger.Errorf("asynq: internal error: both timeout and deadline are not set for the task message: %s", msg.ID)
// 		return p.clock.Now().Add(defaultTimeout)
// 	}
// 	if msg.Timeout != 0 && msg.Deadline != 0 {
// 		deadlineUnix := math.Min(float64(p.clock.Now().Unix()+msg.Timeout), float64(msg.Deadline))
// 		return time.Unix(int64(deadlineUnix), 0)
// 	}
// 	if msg.Timeout != 0 {
// 		return p.clock.Now().Add(time.Duration(msg.Timeout) * time.Second)
// 	}
// 	return time.Unix(msg.Deadline, 0)
// }

// perform calls the handler with the given task.
// If the call returns without panic, it simply returns the value,
// otherwise, it recovers from panic and returns an error.
func (p *processor) perform(ctx context.Context, task *Task) (err error) {
	defer func() {
		if x := recover(); x != nil {
			errMsg := string(debug.Stack())

			p.logger.Errorf("recovering from panic. See the stack trace below for details:\n%s", errMsg)
			_, file, line, ok := runtime.Caller(1) // skip the first frame (panic itself)
			if ok && strings.Contains(file, "runtime/") {
				// The panic came from the runtime, most likely due to incorrect
				// map/slice usage. The parent frame should have the real trigger.
				_, file, line, ok = runtime.Caller(2)
			}

			// Include the file and line number info in the error, if runtime.Caller returned ok.
			if ok {
				err = fmt.Errorf("panic [%s:%d]: %v", file, line, x)
			} else {
				err = fmt.Errorf("panic: %v", x)
			}
			err = &errors.PanicError{
				ErrMsg: errMsg,
			}
		}
	}()
	return p.handler.ProcessTask(ctx, task)
}
