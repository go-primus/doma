package eventbus

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsBus struct {
	nc *nats.Conn
}

type NatsConfig struct {
	Url string
}

func NewNatsBus(cfg NatsConfig) (EventBus, error) {

	url := nats.DefaultURL
	if cfg.Url != "" {
		url = cfg.Url
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsBus{
		nc: conn,
	}, nil
}

func (nb *NatsBus) Request(topic string, data any) {
	b, _ := json.Marshal(data)
	nb.nc.Request(topic, b, time.Second*10)
}

func (nb *NatsBus) Publish(topic string, data any) error {
	b, _ := json.Marshal(data)
	nb.nc.Publish(topic, b)
	return nil
}
func (nb *NatsBus) Subscribe(topic string, fn EventHandler) {
	nb.nc.Subscribe(topic, func(msg *nats.Msg) {
		_msg := Msg(*msg)
		fn(&_msg)
	})
}
func (nb *NatsBus) UnSubscribe(topic string, fn EventHandler) {
	// nb.nc
}
