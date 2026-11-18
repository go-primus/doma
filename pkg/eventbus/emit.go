package eventbus

import (
	"context"
	"encoding/json"
	"errors"
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
	var emitted bool

	// Set up timeout
	timeout := cfg.timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	go func() {
		emitted = b.em.Emit(topic, msg)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("request timeout")
	case <-done:
		if !emitted {
			return nil, nil
		}
		if errMsg := msg.Header.Get(errorHeader); errMsg != "" {
			return nil, errors.New(errMsg)
		}
		return decodeResponse(msg.Data), nil
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
		res, err := fn((*Msg)(msg))
		if msg.Reply == "" {
			return
		}
		if errors.Is(err, ErrNoReply) {
			return
		}

		if err != nil {
			msg.Header.Set(errorHeader, err.Error())
			msg.Data = []byte(err.Error())
			return
		}

		data, mErr := encodeResponse(res)
		if mErr != nil {
			msg.Header.Set(errorHeader, mErr.Error())
			msg.Data = []byte(mErr.Error())
			return
		}
		msg.Data = data
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
