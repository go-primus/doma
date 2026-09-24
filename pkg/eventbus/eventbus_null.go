package eventbus

import (
	"context"
	"sync/atomic"
)

type nullBus struct{}

func NewNullBus() EventBus {
	return &nullBus{}
}

func (b *nullBus) Request(ctx context.Context, topic string, req any, opts ...RequestOption) (any, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, nil
	}
}

func (b *nullBus) Publish(topic string, event any) error {
	return nil
}

func (b *nullBus) Subscribe(topic string, fn EventHandler) (Subscription, error) {
	return &nullSubscription{valid: 1}, nil
}

func (b *nullBus) Unsubscribe(topic string, fn EventHandler) error {
	return nil
}

// nullSubscription implements Subscription
type nullSubscription struct {
	valid int32
}

func (s *nullSubscription) Unsubscribe() error {
	atomic.StoreInt32(&s.valid, 0)
	return nil
}

func (s *nullSubscription) IsValid() bool {
	return atomic.LoadInt32(&s.valid) == 1
}
