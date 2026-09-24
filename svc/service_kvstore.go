package svc

import "encoding/json"

type KvStore interface {
	Put(key string, value any) error
	Delete(key string) error
	DeleteByPrefix(prefix string) error
	Get(key string) (KeyValue, error)
	Count(prefix string) (int64, error)
	List(opts ...ListOption) ([]KeyValue, error)
}

type KeyValue struct {
	Key   string
	Value []byte
}

func (kv *KeyValue) GetString() string {
	if len(kv.Value) == 0 {
		return ""
	}
	var str string
	json.Unmarshal(kv.Value, &str)
	return str
}

func (kv *KeyValue) GetValue(model any) error {
	return json.Unmarshal(kv.Value, model)
}

type ListOption func(*ListOptions)

type ListOptions struct {
	Prefix string
}

func WithListPrefix(prefix string) ListOption {
	return func(lo *ListOptions) {
		lo.Prefix = prefix
	}
}
