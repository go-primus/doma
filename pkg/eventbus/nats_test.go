package eventbus

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"testing"
	"time"
)

// startNatsServer starts a local nats-server for integration tests. The test is
// skipped when the binary is not installed.
func startNatsServer(t *testing.T) string {
	t.Helper()

	bin, err := exec.LookPath("nats-server")
	if err != nil {
		t.Skip("nats-server not installed")
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	cmd := exec.Command(bin, "-a", "127.0.0.1", "-p", fmt.Sprint(port))
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start nats-server: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})

	url := fmt.Sprintf("nats://127.0.0.1:%d", port)
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return url
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("nats-server did not start on %s", url)
	return ""
}

func TestNatsBus_RequestResponse(t *testing.T) {
	url := startNatsServer(t)

	bus, err := NewBus("nats", func(opts Options) { opts["url"] = url })
	if err != nil {
		t.Fatalf("failed to create nats bus: %v", err)
	}

	t.Run("handler result is returned", func(t *testing.T) {
		sub, err := bus.Subscribe("nats.echo", func(msg *Msg) (any, error) {
			return map[string]any{"echo": string(msg.Data)}, nil
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		res, err := bus.Request(context.Background(), "nats.echo", "hello")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		got, ok := res.(map[string]any)
		if !ok || got["echo"] != `"hello"` {
			t.Fatalf("unexpected response: %#v", res)
		}
	})

	t.Run("handler error is returned", func(t *testing.T) {
		sub, err := bus.Subscribe("nats.error", func(msg *Msg) (any, error) {
			return nil, fmt.Errorf("handler failed")
		})
		if err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
		defer sub.Unsubscribe()

		_, err = bus.Request(context.Background(), "nats.error", "hello")
		if err == nil || err.Error() != "handler failed" {
			t.Fatalf("expected handler error, got %v", err)
		}
	})
}
