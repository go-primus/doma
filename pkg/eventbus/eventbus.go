package eventbus

import (
	"context"
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

// EventHandler processes received messages
type EventHandler func(msg *Msg)

// Msg wraps NATS message
type Msg = nats.Msg
