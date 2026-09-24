package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// TestEmitBus_Subscribe tests EmitBus Subscribe functionality
func TestEmitBus_Subscribe(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("subscribe returns subscription", func(t *testing.T) {
		sub, err := bus.Subscribe("test.topic", func(msg *Msg) {
			slog.Info("receive", "msg", string(msg.Data))
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if sub == nil {
			t.Error("expected non-nil subscription")
		}
	})

	t.Run("multiple subscriptions to same topic", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			_, err := bus.Subscribe("multi.topic", func(msg *Msg) {
				slog.Info("receive", "msg", string(msg.Data))
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}

		_ = bus.Publish("multi.topic", "hello world")
	})
}

// TestEmitBus_Publish tests EmitBus Publish functionality
func TestEmitBus_Publish(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("publish returns no error", func(t *testing.T) {
		err := bus.Publish("test.topic", "data")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestEmitBus_Request tests EmitBus Request functionality
func TestEmitBus_Request(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("request with handler", func(t *testing.T) {
		sub, err := bus.Subscribe("request.topic", func(msg *Msg) {
			msg.Data = []byte("response data")
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		ctx := context.Background()
		data, err := bus.Request(ctx, "request.topic", "request data")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if data == nil {
			t.Error("expected non-nil data")
		}
	})
}

// TestWithTimeout tests RequestOption WithTimeout
func TestWithTimeout(t *testing.T) {
	t.Run("WithTimeout sets timeout", func(t *testing.T) {
		opt := WithTimeout(10 * time.Second)
		cfg := &requestConfig{}
		opt(cfg)
		if cfg.timeout != 10*time.Second {
			t.Errorf("expected timeout 10s, got %v", cfg.timeout)
		}
	})
}

// TestEmitBus_MessageDelivery tests that published messages are received by subscribers
func TestEmitBus_MessageDelivery(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("subscriber receives published message", func(t *testing.T) {
		var receivedMu sync.Mutex
		var receivedData []byte

		sub, err := bus.Subscribe("delivery.test", func(msg *Msg) {
			receivedMu.Lock()
			receivedData = msg.Data
			receivedMu.Unlock()
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		testData := map[string]any{"message": "hello, world!"}
		err = bus.Publish("delivery.test", testData)
		if err != nil {
			t.Fatalf("failed to publish: %v", err)
		}

		time.Sleep(50 * time.Millisecond)

		receivedMu.Lock()
		defer receivedMu.Unlock()
		if receivedData == nil {
			t.Error("subscriber did not receive message")
		}
		// Data is JSON marshaled
		if len(receivedData) == 0 {
			t.Error("received data is empty")
		}
	})

	t.Run("multiple subscribers receive same message", func(t *testing.T) {
		var mu sync.Mutex
		receivedCount := 0
		expectedCount := 3

		for i := 0; i < expectedCount; i++ {
			_, err := bus.Subscribe("multi.delivery", func(msg *Msg) {
				mu.Lock()
				receivedCount++
				mu.Unlock()
			})
			if err != nil {
				t.Fatalf("failed to subscribe: %v", err)
			}
		}

		err = bus.Publish("multi.delivery", "broadcast message")
		if err != nil {
			t.Fatalf("failed to publish: %v", err)
		}

		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		if receivedCount != expectedCount {
			t.Errorf("expected %d subscribers to receive message, got %d", expectedCount, receivedCount)
		}
	})

	t.Run("subscriber receives correct topic", func(t *testing.T) {
		var receivedTopic string

		sub, err := bus.Subscribe("topic.test", func(msg *Msg) {
			receivedTopic = msg.Subject
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		err = bus.Publish("topic.test", "data")
		if err != nil {
			t.Fatalf("failed to publish: %v", err)
		}

		time.Sleep(50 * time.Millisecond)

		if receivedTopic != "topic.test" {
			t.Errorf("expected topic %q, got %q", "topic.test", receivedTopic)
		}
	})

	t.Run("struct data delivered correctly", func(t *testing.T) {
		type TestData struct {
			Name string
			Age  int
		}
		var receivedData []byte

		sub, err := bus.Subscribe("struct.test", func(msg *Msg) {
			receivedData = msg.Data
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		testData := TestData{Name: "john", Age: 30}
		err = bus.Publish("struct.test", testData)
		if err != nil {
			t.Fatalf("failed to publish: %v", err)
		}

		time.Sleep(50 * time.Millisecond)

		if receivedData == nil {
			t.Error("subscriber did not receive struct data")
		}
	})
}

// TestEmitBus_RequestResponse tests request-response pattern with message delivery
func TestEmitBus_RequestResponse(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("handler receives request and responds", func(t *testing.T) {
		sub, err := bus.Subscribe("echo.service", func(msg *Msg) {
			msg.Data = []byte("echo response")
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		ctx := context.Background()
		data, err := bus.Request(ctx, "echo.service", "hello")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}

		if data == nil {
			t.Fatal("expected response data, got nil")
		}
	})

	t.Run("request with timeout returns error", func(t *testing.T) {
		sub, err := bus.Subscribe("slow.service", func(msg *Msg) {
			time.Sleep(2 * time.Second)
			msg.Data = []byte("slow response")
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		_, err = bus.Request(ctx, "slow.service", "test")
		if err == nil {
			t.Error("expected timeout error, got nil")
		}
	})

	t.Run("multiple requests handled correctly", func(t *testing.T) {
		requestCount := 0
		sub, err := bus.Subscribe("counter.service", func(msg *Msg) {
			requestCount++
			msg.Data = []byte(fmt.Sprintf("response %d", requestCount))
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		ctx := context.Background()
		for i := 0; i < 5; i++ {
			data, err := bus.Request(ctx, "counter.service", "test")
			if err != nil {
				t.Fatalf("request %d failed: %v", i, err)
			}
			if data == nil {
				t.Errorf("request %d: expected response, got nil", i)
			}
		}
	})
}

// TestEmitBus_UnsubscribeDelivery tests message delivery after unsubscribe
func TestEmitBus_UnsubscribeDelivery(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("subscription validity check", func(t *testing.T) {
		sub, err := bus.Subscribe("validity.test", func(msg *Msg) {})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}

		if !sub.IsValid() {
			t.Error("new subscription should be valid")
		}

		err = sub.Unsubscribe()
		if err != nil {
			t.Fatalf("failed to unsubscribe: %v", err)
		}

		if sub.IsValid() {
			t.Error("unsubscribed subscription should not be valid")
		}
	})

	t.Run("subscription count tracking", func(t *testing.T) {
		var receivedCount int

		sub, err := bus.Subscribe("counting.test", func(msg *Msg) {
			receivedCount++
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}

		_ = bus.Publish("counting.test", "first")
		time.Sleep(20 * time.Millisecond)

		err = sub.Unsubscribe()
		if err != nil {
			t.Fatalf("failed to unsubscribe: %v", err)
		}

		_ = bus.Publish("counting.test", "second")
		time.Sleep(20 * time.Millisecond)

		if receivedCount < 1 {
			t.Error("should have received at least one message before unsubscribe")
		}
	})
}

// TestEmitBus_TopicDelivery tests various topic delivery scenarios
func TestEmitBus_TopicDelivery(t *testing.T) {
	bus, err := NewBus("emit")
	if err != nil {
		t.Fatalf("failed to create emit bus: %v", err)
	}

	t.Run("exact topic delivery", func(t *testing.T) {
		var receivedTopics []string

		sub1, err := bus.Subscribe("exact.topic", func(msg *Msg) {
			receivedTopics = append(receivedTopics, msg.Subject)
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub1.Unsubscribe()

		_ = bus.Publish("exact.topic", "data1")
		_ = bus.Publish("other.topic", "data2")
		_ = bus.Publish("exact.topic.sub", "data3")

		time.Sleep(50 * time.Millisecond)

		foundExact := false
		foundOther := false
		foundSub := false
		for _, topic := range receivedTopics {
			if topic == "exact.topic" {
				foundExact = true
			}
			if topic == "other.topic" {
				foundOther = true
			}
			if topic == "exact.topic.sub" {
				foundSub = true
			}
		}

		if !foundExact {
			t.Error("should receive exact topic messages")
		}
		if foundOther {
			t.Error("should not receive other topic messages")
		}
		if foundSub {
			t.Error("should not receive subtopic messages")
		}
	})

	t.Run("multiple topic subscriptions", func(t *testing.T) {
		var receivedOnA, receivedOnB bool

		sub, err := bus.Subscribe("topic.a", func(msg *Msg) {
			receivedOnA = true
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		sub2, err := bus.Subscribe("topic.b", func(msg *Msg) {
			receivedOnB = true
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub2.Unsubscribe()

		_ = bus.Publish("topic.a", "data1")
		_ = bus.Publish("topic.b", "data2")

		time.Sleep(50 * time.Millisecond)

		if !receivedOnA {
			t.Error("should receive topic.a message")
		}
		if !receivedOnB {
			t.Error("should receive topic.b message")
		}
	})
}

// BenchmarkNullBus_Publish benchmarks Publish operation
func BenchmarkNullBus_Publish(b *testing.B) {
	bus := NewNullBus()
	topic := "bench.topic"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bus.Publish(topic, map[string]any{"data": i})
	}
}

// BenchmarkNullBus_Request benchmarks Request operation
func BenchmarkNullBus_Request(b *testing.B) {
	bus := NewNullBus()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bus.Request(ctx, "bench.topic", map[string]any{"data": i})
	}
}

// BenchmarkNullBus_Subscribe benchmarks Subscribe operation
func BenchmarkNullBus_Subscribe(b *testing.B) {
	bus := NewNullBus()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sub, _ := bus.Subscribe(fmt.Sprintf("bench.topic.%d", i), func(msg *Msg) {})
		_ = sub.Unsubscribe()
	}
}

// BenchmarkEmitBus_Publish benchmarks EmitBus Publish operation
func BenchmarkEmitBus_Publish(b *testing.B) {
	bus, _ := NewBus("emit")
	topic := "bench.topic"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bus.Publish(topic, map[string]any{"data": i})
	}
}

// BenchmarkEmitBus_Subscribe benchmarks EmitBus Subscribe operation
func BenchmarkEmitBus_Subscribe(b *testing.B) {
	bus, _ := NewBus("emit")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sub, _ := bus.Subscribe(fmt.Sprintf("bench.topic.%d", i), func(msg *Msg) {})
		_ = sub.Unsubscribe()
	}
}
