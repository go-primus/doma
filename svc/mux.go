package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/go-primus/doma/core/common/event"
)

type Handler interface {
	HandleEvent(context.Context, event.Event) error
}

type HandlerFunc func(context.Context, event.Event) error

func (fn HandlerFunc) HandleEvent(ctx context.Context, event event.Event) error {
	return fn(ctx, event)
}

// ////////
type EventMux struct {
	s     Service
	_name string
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry
	mws   []MiddlewareFunc
}

type muxEntry struct {
	h       Handler
	pattern string
}

type MiddlewareFunc func(Handler) Handler

// //

func NewEventMux(muxName string) *EventMux {
	mux := new(EventMux)
	mux._name = muxName
	return mux
}

func NewServiceEventMux(s Service) *EventMux {
	mux := NewEventMux(s.Name())
	mux.s = s
	return mux
}

func (mux *EventMux) name() string {

	if mux._name != "" {
		return mux._name
	}

	if mux.s == nil {
		return ""
	}
	return mux.s.Name()
}

func (mux *EventMux) Subscribe(topics ...string) {
	s := mux.s

	if s == nil {
		return
	}
	for _, topic := range topics {
		s.Subscribe(topic, mux.OnEvent)
	}
}

func (mux *EventMux) Use(mws ...MiddlewareFunc) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	mux.mws = append(mux.mws, mws...)
}

func (mux *EventMux) OnEvent(ctx context.Context, topic string, event event.Event) error {

	// 解析 event
	eventType := event.EventType
	if eventType == "" {
		return nil
	}

	h, pattern := mux.Handler(string(eventType), event)
	if pattern == "" {
		return nil
	}
	slog.Debug("=======[" + mux.name() + "] ---- handle event: " + string(event.EventType) + " ==============")
	err := h.HandleEvent(ctx, event)
	if err != nil {
		slog.Error("=======["+mux.name()+"] ---- handle event end:", "err", err)
	} else {
		slog.Debug("=======[" + mux.name() + "] ---- handle event end: ok ==============")
	}
	return nil
}

func (mux *EventMux) Handler(topic string, event event.Event) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, pattern = mux.match(topic)
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	for i := len(mux.mws) - 1; i >= 0; i-- {
		h = mux.mws[i](h)
	}
	return h, pattern
}

// Find a handler on a handler map given a typename string.
// Most-specific (longest) pattern wins.
func (mux *EventMux) match(typename string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[typename]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.
	// mux.es contains all patterns from longest to shortest.
	for _, e := range mux.es {
		if strings.HasPrefix(typename, e.pattern) {
			return e.h, e.pattern
		}
	}
	return nil, ""

}

func (mux *EventMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if strings.TrimSpace(pattern) == "" {
		return
	}

	if handler == nil {
		return
	}

	if _, exist := mux.m[pattern]; exist {
		slog.Warn("event mux: multiple registrations for " + pattern)
		return
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	e := muxEntry{
		h:       handler,
		pattern: pattern,
	}

	mux.m[pattern] = e
	mux.es = append(mux.es, e)

}

func (mux *EventMux) HandleFunc(pattern string, handler HandlerFunc) {
	if handler == nil {
		return
	}
	mux.Handle(pattern, handler)
}

/////

func NotFound(ctx context.Context, event event.Event) error {
	return nil
}

func NotFoundHandler() Handler {
	return HandlerFunc(NotFound)
}

func WrapHandle[T any](fn func(ctx context.Context, eventType string, payload T) error) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, e event.Event) error {

		bb, _ := json.Marshal(e.Payload)

		payload := new(T)
		err := json.Unmarshal(bb, payload)
		if err != nil {
			fmt.Println("unmarshal err:", err)
			return err
		}

		return fn(ctx, string(e.EventType), *payload)

	})
}
