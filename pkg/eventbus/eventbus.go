package eventbus

import (
	"github.com/nats-io/nats.go"
)

//https://libuba.com/2020/11/02/golang%E5%8C%85%E5%BE%AA%E7%8E%AF%E5%BC%95%E7%94%A8%E7%9A%84%E5%87%A0%E7%A7%8D%E8%A7%A3%E5%86%B3%E6%96%B9%E6%A1%88/

type Msg nats.Msg

type Event struct{}
type EventHandler func(msg *Msg)

type EventBus interface {
	// Request(topic string, req any, duration time.Duration) (any, error)
	Publish(topc string, data any) error
	Subscribe(topic string, fn EventHandler)
	UnSubscribe(topic string, fn EventHandler)
}
