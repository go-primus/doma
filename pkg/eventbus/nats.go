package eventbus

import (
	"context"
	"encoding/json"
	"errors"
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

	if ctx == nil {
		return nil, nats.ErrInvalidContext
	}

	timeout := cfg.timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	// An explicit deadline on ctx takes priority; otherwise apply the
	// request timeout (from WithTimeout or the default).
	reqCtx := ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	b, _ := json.Marshal(req)
	res, err := nb.nc.RequestWithContext(reqCtx, topic, b)
	if err != nil {
		return nil, err
	}

	if errMsg := res.Header.Get(errorHeader); errMsg != "" {
		return nil, errors.New(errMsg)
	}

	return decodeResponse(res.Data), nil
}

func (nb *NatsBus) Publish(topic string, data any) error {
	b, _ := json.Marshal(data)
	return nb.nc.Publish(topic, b)
}

func (nb *NatsBus) Subscribe(topic string, fn EventHandler) (Subscription, error) {
	sub, err := nb.nc.Subscribe(topic, func(msg *nats.Msg) {
		res, err := fn(msg)
		if msg.Reply == "" {
			return
		}
		if errors.Is(err, ErrNoReply) {
			return
		}

		if err != nil {
			_ = respondError(msg, err)
			return
		}

		data, mErr := encodeResponse(res)
		if mErr != nil {
			_ = respondError(msg, mErr)
			return
		}
		_ = msg.Respond(data)
	})
	if err != nil {
		return nil, err
	}
	return &natsSubscription{sub: sub}, nil
}

// respondError publishes an error response. The error is carried in a header
// so Request can distinguish it from a normal payload.
func respondError(msg *nats.Msg, err error) error {
	resp := nats.NewMsg(msg.Reply)
	resp.Header.Set(errorHeader, err.Error())
	resp.Data = []byte(err.Error())
	return msg.RespondMsg(resp)
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
