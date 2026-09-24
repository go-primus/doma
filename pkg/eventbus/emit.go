package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jiyeyuran/go-eventemitter"
	"github.com/nats-io/nats.go"
)

type emitDriver struct{}

func (d *emitDriver) Name() string {
	return "emit"
}

func (d *emitDriver) New(...Option) (EventBus, error) {

	em := eventemitter.New(eventemitter.WithQueueSize(10000))
	return &emitBus{
		em:   em,
		subs: make(map[string]any),
	}, nil
}

type emitBus struct {
	em   eventemitter.IEventEmitter
	subs map[string]any
}

func (b *emitBus) Request(ctx context.Context, topic string, req any, opts ...RequestOption) (any, error) {
	cfg := &requestConfig{timeout: 5 * time.Second}
	for _, opt := range opts {
		opt(cfg)
	}

	bb, _ := json.Marshal(req)

	msg := nats.NewMsg(topic)
	msg.Reply = uuid.NewString()
	msg.Data = bb

	done := make(chan struct{})
	var result any

	// Set up timeout
	timeout := cfg.timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	go func() {
		ok := b.em.Emit(topic, msg)
		if !ok {
			result = nil
		} else {
			result = msg.Data
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("request timeout")
	case <-done:
		return result, nil
	}
}

func (b *emitBus) Publish(topic string, data any) error {
	bb, _ := json.Marshal(data)

	msg := nats.NewMsg(topic)
	msg.Data = bb

	b.em.AsyncEmit(topic, msg)
	return nil
}

func (b *emitBus) Subscribe(topic string, fn EventHandler) (Subscription, error) {
	token := b.em.On(topic, func(msg *nats.Msg) {
		fn((*Msg)(msg))
	})
	sub := &emitSubscription{topic: topic, token: token, bus: b}
	b.subs[topic] = token
	return sub, nil
}

func (b *emitBus) Unsubscribe(topic string, fn EventHandler) error {
	if token, ok := b.subs[topic]; ok {
		b.em.Off(topic, token)
		delete(b.subs, topic)
	}
	return nil
}

// emitSubscription implements Subscription
type emitSubscription struct {
	topic string
	token any
	bus   *emitBus
}

func (s *emitSubscription) Unsubscribe() error {
	delete(s.bus.subs, s.topic)
	return nil
}

func (s *emitSubscription) IsValid() bool {
	_, ok := s.bus.subs[s.topic]
	return ok
}
