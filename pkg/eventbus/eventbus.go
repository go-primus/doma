package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

// EventBus provides event publish/subscribe capabilities
type EventBus interface {
	// Request sends a request and waits for a response
	Request(ctx context.Context, topic string, req any, opts ...RequestOption) (any, error)
	// Publish publishes an event to a topic
	Publish(topic string, data any) error
	// Subscribe subscribes to a topic, returns Subscription for lifecycle management
	Subscribe(topic string, fn EventHandler) (Subscription, error)
	// Unsubscribe cancels a subscription (kept for compatibility, prefer Subscription.Unsubscribe)
	Unsubscribe(topic string, fn EventHandler) error
}

// Subscription represents an active subscription
type Subscription interface {
	// Unsubscribe cancels the subscription
	Unsubscribe() error
	// IsValid checks if the subscription is still valid
	IsValid() bool
}

// RequestOption configures request behavior
type RequestOption func(*requestConfig)

type requestConfig struct {
	timeout time.Duration
	headers map[string]string
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) RequestOption {
	return func(c *requestConfig) { c.timeout = timeout }
}

// errorHeader is used to carry handler errors back to Request callers.
const errorHeader = "Doma-Error"

// ErrNoReply can be returned by an EventHandler to indicate that the response
// will be sent asynchronously by the handler itself. The bus then skips the
// automatic reply.
var ErrNoReply = errors.New("eventbus: no reply")

// EventHandler processes received messages.
//
// When the message has a reply subject (i.e. it was sent by Request), the
// returned value is encoded and sent back as the response. A non-nil error is
// sent back as an error response and surfaces as the Request error.
type EventHandler func(msg *Msg) (any, error)

// Msg wraps NATS message
type Msg = nats.Msg

// encodeResponse converts a handler result into the response payload.
func encodeResponse(res any) ([]byte, error) {
	if res == nil {
		return nil, nil
	}
	if b, ok := res.([]byte); ok {
		return b, nil
	}
	return json.Marshal(res)
}

// decodeResponse converts a response payload into a Go value. Non-JSON payloads
// are returned as raw bytes.
func decodeResponse(data []byte) any {
	if len(data) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return data
	}
	return v
}
