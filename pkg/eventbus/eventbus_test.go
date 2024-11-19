package eventbus

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestEventBus(t *testing.T) {

	bus, _ := NewEmitBus()
	bus.Subscribe("file.create", func(msg *nats.Msg) {

		data := msg.Data
		fmt.Println("---recive event:", msg.Subject, "----", string(data))

	})

	bus.Publish("file.create", "create file")
}
