package broker

import (
	"context"
	"time"

	"github.com/jiyeyuran/go-eventemitter"
	"github.com/primus/primus/pkg/asynctask/internal/base"
)

type EmitBroker struct {
	emitter eventemitter.IEventEmitter
}

// Dequeue implements base.Broker.
func (e *EmitBroker) Dequeue(qnames ...string) (*base.TaskMessage, time.Time, error) {
	panic("unimplemented")
	// e.emitter.On()
}

// Enqueue implements base.Broker.
func (e *EmitBroker) Enqueue(ctx context.Context, msg *base.TaskMessage) error {
	panic("unimplemented")
}

// EnqueueUnique implements base.Broker.
func (e *EmitBroker) EnqueueUnique(ctx context.Context, msg *base.TaskMessage, ttl time.Duration) error {
	panic("unimplemented")
}

func NewEmitBroker() base.Broker {
	return &EmitBroker{
		emitter: eventemitter.NewEventEmitter(),
	}
}
