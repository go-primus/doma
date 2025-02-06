package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-primus/doma/core/common/event"
)

func makeFakeHandler() Handler {
	return HandlerFunc(
		func(ctx context.Context, e event.Event) {
			fmt.Println("----handle event----", e)
		},
	)
}

func wrapHandle[T any](fn func(eventType string, payload T)) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, e event.Event) {

		bb, _ := json.Marshal(e.Payload)

		payload := new(T)
		err := json.Unmarshal(bb, payload)
		if err != nil {
			fmt.Println("unmarshal err:", err)
			return
		}

		fn(string(e.EventType), *payload)

	})
}

func TestMux(t *testing.T) {

	mux := NewEventMux()
	mux.Handle("xxxx", makeFakeHandler())
	mux.HandleFunc("yyyy", wrapHandle(func(typ string, payload struct {
		Name string
		Size int64
	}) {
		fmt.Println("---", typ, "----", payload)
	}))
	mux.OnEvent(context.Background(), "xxxx", event.NewEvent("xxxx", struct {
		Name string
		Size int64
	}{
		Name: "nameooo",
		Size: 2348,
	}))

}
