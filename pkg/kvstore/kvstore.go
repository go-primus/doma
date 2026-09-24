package kvstore

import (
	"encoding/json"
	"errors"
)

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrBucketNotFound = errors.New("bucket not found")
)

type KeyValue struct {
	Key   string
	Value []byte
}

func (kv *KeyValue) GetString() string {

	var str string
	json.Unmarshal(kv.Value, &str)
	return str
}

func (kv *KeyValue) GetValue(model any) error {
	return json.Unmarshal(kv.Value, model)
}

type KvStore interface {
	Put(bucket, key string, value []byte) error
	Delete(bucket, key string) error
	DeleteByPrefix(bucket, prefix string) error
	Get(bucket, key string) (KeyValue, error)
	Count(bucket string, prefix string) (int64, error)
	List(bucket string, opts ...ListOption) ([]KeyValue, error)
	Close() error
}

type ListOption func(*ListOptions)

type ListOptions struct {
	Prefix  string
	Limit   int
	Offset  int
	Reverse bool
}

func WithPrefix(prefix string) ListOption {
	return func(o *ListOptions) {
		o.Prefix = prefix
	}
}

func WithLimit(limit int) ListOption {
	return func(o *ListOptions) {
		o.Limit = limit
	}
}

func WithOffset(offset int) ListOption {
	return func(o *ListOptions) {
		o.Offset = offset
	}
}

func WithReverse() ListOption {
	return func(o *ListOptions) {
		o.Reverse = true
	}
}
