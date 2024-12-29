package asynctask

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-primus/doma/pkg/asynctask/internal/base"
)

func (p *processor) handleSucceededMessage(msg *base.TaskMessage) {
	if msg.Retention > 0 {
		p.markAsComplete(msg)
	} else {
		p.markAsDone(msg)
	}
}

func (p *processor) markAsComplete(msg *base.TaskMessage) {

	p.logger.Info("---complete----")
	// ctx := context.Background()
	// err := p.broker.MarkAsComplete(ctx, msg)

	// if err != nil {
	// 	errMsg := fmt.Sprintf("Could not move task id=%s type=%q from %q to %q:  %+v",
	// 		msg.ID, msg.Type, base.ActiveKey(msg.Queue), base.CompletedKey(msg.Queue), err)
	// 	p.logger.Warnf("%s; Will retry syncing", errMsg)
	// 	p.syncRequestCh <- &syncRequest{
	// 		fn: func() error {
	// 			return p.broker.MarkAsComplete(ctx, msg)
	// 		},
	// 		errMsg:   errMsg,
	// 		deadline: l.Deadline(),
	// 	}
	// }
}

func (p *processor) markAsDone(msg *base.TaskMessage) {
	fmt.Println("----done----")
	// ctx, _ := context.WithDeadline(context.Background(), l.Deadline())
	// err := p.broker.Done(ctx, msg)
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Could not remove task id=%s type=%q from %q err: %+v", msg.ID, msg.Type, base.ActiveKey(msg.Queue), err)
	// 	p.logger.Warnf("%s; Will retry syncing", errMsg)
	// 	p.syncRequestCh <- &syncRequest{
	// 		fn: func() error {
	// 			return p.broker.Done(ctx, msg)
	// 		},
	// 		errMsg:   errMsg,
	// 		deadline: l.Deadline(),
	// 	}
	// }
}

// SkipRetry is used as a return value from Handler.ProcessTask to indicate that
// the task should not be retried and should be archived instead.
var SkipRetry = errors.New("skip retry for the task")

func (p *processor) handleFailedMessage(ctx context.Context, msg *base.TaskMessage, err error) {
	p.logger.Info("--handle failed message---err:", err)
	// if p.errHandler != nil {
	// 	p.errHandler.HandleError(ctx, NewTask(msg.Type, msg.Payload), err)
	// }
	// if !p.isFailureFunc(err) {
	// 	// retry the task without marking it as failed
	// 	p.retry(l, msg, err, false /*isFailure*/)
	// 	return
	// }
	// if msg.Retried >= msg.Retry || errors.Is(err, SkipRetry) {
	// 	p.logger.Warnf("Retry exhausted for task id=%s", msg.ID)
	// 	p.archive(l, msg, err)
	// } else {
	// 	p.retry(l, msg, err, true /*isFailure*/)
	// }
}
