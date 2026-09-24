package store

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestMemoryStore_PutAndGet(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	err := s.Put(ctx, "order:1001:status", "paid")
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	kv, err := s.Get(ctx, "order:1001:status")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if kv.Key != "order:1001:status" {
		t.Errorf("expected key 'order:1001:status', got %q", kv.Key)
	}
	if kv.Value != "paid" {
		t.Errorf("expected value 'paid', got %v", kv.Value)
	}
}

func TestMemoryStore_Get_NotFound(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	_, err := s.Get(ctx, "nonexistent")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryStore_Exists(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	exists, err := s.Exists(ctx, "order:1001:status")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected false for nonexistent key")
	}

	s.Put(ctx, "order:1001:status", "paid")

	exists, err = s.Exists(ctx, "order:1001:status")
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected true for existing key")
	}
}

func TestMemoryStore_Update(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	s.Put(ctx, "order:1001:status", "pending")

	err := s.Update(ctx, "order:1001:status", "paid")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	kv, _ := s.Get(ctx, "order:1001:status")
	if kv.Value != "paid" {
		t.Errorf("expected 'paid', got %v", kv.Value)
	}
}

func TestMemoryStore_Update_NotFound(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	err := s.Update(ctx, "nonexistent", "value")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	s.Put(ctx, "order:1001:status", "paid")

	err := s.Delete(ctx, "order:1001:status")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = s.Get(ctx, "order:1001:status")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Error("expected ErrKeyNotFound after delete")
	}
}

func TestMemoryStore_Delete_NotFound(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	err := s.Delete(ctx, "nonexistent")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryStore_List(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	keys := []string{
		"order:1001:status",
		"order:1001:result",
		"order:1002:status",
		"task:3001:progress",
		"task:3001:result",
	}
	for _, k := range keys {
		s.Put(ctx, k, "value")
	}

	result, err := s.List(ctx, ListOptions{Prefix: ""})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Items) != len(keys) {
		t.Errorf("expected %d items, got %d", len(keys), len(result.Items))
	}

	result, err = s.List(ctx, ListOptions{Prefix: "order:"})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Items) != 3 {
		t.Errorf("expected 3 order items, got %d", len(result.Items))
	}
}

func TestMemoryStore_List_Limit(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	s.Put(ctx, "task:3001:result", "value1")
	s.Put(ctx, "task:3002:result", "value2")
	s.Put(ctx, "task:3003:result", "value3")
	s.Put(ctx, "task:3004:result", "value4")
	s.Put(ctx, "task:3005:result", "value5")

	result, err := s.List(ctx, ListOptions{Prefix: "task:", Limit: 3})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.Items))
	}
}

func TestMemoryStore_Watch_Create(t *testing.T) {
	s := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := s.Watch(ctx, "order:")
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}

	go func() {
		s.Put(ctx, "order:1001:status", "paid")
	}()

	select {
	case event := <-ch:
		if event.Type != EventCreate {
			t.Errorf("expected EventCreate, got %v", event.Type)
		}
		if event.Key != "order:1001:status" {
			t.Errorf("expected key 'order:1001:status', got %q", event.Key)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestMemoryStore_Watch_Update(t *testing.T) {
	s := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Put(ctx, "order:1001:status", "pending")

	ch, _ := s.Watch(ctx, "order:")

	go func() {
		s.Update(ctx, "order:1001:status", "paid")
	}()

	select {
	case event := <-ch:
		if event.Type != EventUpdate {
			t.Errorf("expected EventUpdate, got %v", event.Type)
		}
		if event.Value != "paid" {
			t.Errorf("expected value 'paid', got %v", event.Value)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestMemoryStore_Watch_Delete(t *testing.T) {
	s := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Put(ctx, "order:1001:status", "paid")

	ch, _ := s.Watch(ctx, "order:")

	go func() {
		s.Delete(ctx, "order:1001:status")
	}()

	select {
	case event := <-ch:
		if event.Type != EventDelete {
			t.Errorf("expected EventDelete, got %v", event.Type)
		}
		if event.Value != nil {
			t.Errorf("expected nil value on delete, got %v", event.Value)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestMemoryStore_Watch_MultipleWatchers(t *testing.T) {
	s := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1, _ := s.Watch(ctx, "order:")
	ch2, _ := s.Watch(ctx, "order:1001:")

	go func() {
		s.Put(ctx, "order:1001:status", "paid")
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case <-ch1:
		case <-time.After(time.Second):
			t.Error("timeout on watcher 1")
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-ch2:
		case <-time.After(time.Second):
			t.Error("timeout on watcher 2")
		}
	}()

	wg.Wait()
}
