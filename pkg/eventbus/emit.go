package eventbus

import (
	"encoding/json"

	"github.com/jiyeyuran/go-eventemitter"
	"github.com/nats-io/nats.go"
)

type emitBus struct {
	em eventemitter.IEventEmitter
}

func (b *emitBus) Request(topic string, data any) error {
	bb, _ := json.Marshal(data)
	b.em.Emit(topic, bb)

	return nil
}

// Publish implements EventBus.
func (b *emitBus) Publish(topc string, data any) error {
	bb, _ := json.Marshal(data)

	b.em.SafeEmit(topc, bb)
	return nil
}

// Subscribe implements EventBus.
func (b *emitBus) Subscribe(topic string, fn EventHandler) {

	b.em.On(topic, func(data []byte) {
		msg := nats.Msg{}
		msg.Data = data
		msg.Subject = topic
		fn(&msg)
	})
}

// UnSubscribe implements EventBus.
func (*emitBus) UnSubscribe(topic string, fn EventHandler) {
	panic("unimplemented")
}

func NewEmitBus() (EventBus, error) {

	//
	em := eventemitter.NewEventEmitter(eventemitter.WithQueueSize(10000))
	return &emitBus{
		em: em,
	}, nil
}
