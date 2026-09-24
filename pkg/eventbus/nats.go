package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

type natsDriver struct{}

func (d *natsDriver) Name() string {
	return "nats"
}

func (d *natsDriver) New(opts ...Option) (EventBus, error) {

	optx := Options{}
	for _, opt := range opts {
		opt(optx)
	}

	url := optx.GetString("url")
	if url == "" {
		url = nats.DefaultURL
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsBus{
		nc: conn,
	}, nil
}

type NatsBus struct {
	nc *nats.Conn
}

type NatsConfig struct {
	Url string
}

func (nb *NatsBus) Request(ctx context.Context, topic string, req any, opts ...RequestOption) (any, error) {
	cfg := &requestConfig{timeout: 5 * time.Second}
	for _, opt := range opts {
		opt(cfg)
	}

	b, _ := json.Marshal(req)
	res, err := nb.nc.RequestWithContext(ctx, topic, b)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (nb *NatsBus) Publish(topic string, data any) error {
	b, _ := json.Marshal(data)
	return nb.nc.Publish(topic, b)
}

func (nb *NatsBus) Subscribe(topic string, fn EventHandler) (Subscription, error) {
	sub, err := nb.nc.Subscribe(topic, func(msg *nats.Msg) {
		_msg := Msg(*msg)
		fn(&_msg)
	})
	if err != nil {
		return nil, err
	}
	return &natsSubscription{sub: sub}, nil
}

func (nb *NatsBus) Unsubscribe(topic string, fn EventHandler) error {
	return nil
}

// natsSubscription implements Subscription
type natsSubscription struct {
	sub *nats.Subscription
}

func (s *natsSubscription) Unsubscribe() error {
	return s.sub.Unsubscribe()
}

func (s *natsSubscription) IsValid() bool {
	return s.sub.IsValid()
}
