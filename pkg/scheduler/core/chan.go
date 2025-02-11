package core

import (
	"encoding/json"
	"fmt"

	"github.com/go-primus/doma/pkg/eventbus"
)

type ChannelBuilder struct{}

//////////////////

type Channel[T any] interface {
	// Send(T)
	// Recv() <-chan T
	ChanSender[T]
	ChanRecver[T]
}

type ChanSender[T any] interface {
	Send(T)
}

type ChanRecver[T any] interface {
	Recv() <-chan T
}

/////////////////////////

type gochan[T any] struct {
	msgchan chan T
}

func (ch *gochan[T]) Send(msg T) {
	ch.msgchan <- msg
}

func (ch *gochan[T]) Recv() <-chan T {

	return ch.msgchan
}

func demo() {
	c := &gochan[string]{}

	ch := c.Recv()
	msg := <-ch

	fmt.Println("msg:", msg)
}

type busChan[T any] struct {
	topic    string
	bus      eventbus.EventBus
	dispatch chan T
}

func (s *busChan[T]) Send(msg T) {
	s.bus.Publish(s.topic, msg)
}

func (s *busChan[T]) Recv() <-chan T {

	s.bus.Subscribe(s.topic, func(msg *eventbus.Msg) {

		var item T
		json.Unmarshal(msg.Data, &item)

		s.dispatch <- item
	})

	return s.dispatch
}
