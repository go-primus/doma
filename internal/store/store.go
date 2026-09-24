package store

import (
	"context"
	"errors"
)

// Predefined errors
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key already exists")
)

// EventType represents the type of watch event
type EventType int

const (
	EventCreate EventType = iota
	EventUpdate
	EventDelete
)

// WatchEvent represents a single watch event
type WatchEvent struct {
	Type EventType
	KeyValue
}

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value any
}

// Store is a key-value store with observability support
type Store interface {
	// Basic operations
	Put(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (KeyValue, error)
	Delete(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Update updates existing key, returns ErrKeyNotFound if not exists
	Update(ctx context.Context, key string, value any) error

	// List returns all keys, use cursor for pagination
	List(ctx context.Context, opts ListOptions) (ListResult, error)

	// Watch subscribes to key changes
	Watch(ctx context.Context, prefix string) (<-chan WatchEvent, error)
}

// ListOptions for pagination
type ListOptions struct {
	Prefix string // filter by prefix
	Limit  int    // max results (default 100, max 1000)
	Cursor string // pagination cursor from previous response
}

// ListResult contains paginated results
type ListResult struct {
	Items  []KeyValue
	Cursor string // empty means no more results
}
