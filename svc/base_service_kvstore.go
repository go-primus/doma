package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-primus/doma/pkg/kvstore"
)

type BaseServiceKvStore struct {
	kvstore kvstore.KvStore
	s       Service
}

// Count implements [ProcessWithKey].
func (s *BaseServiceKvStore) Count(prefix string) (int64, error) {
	if s.kvstore == nil {
		return 0, fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}

	return s.kvstore.Count(s.s.Name(), prefix)
}

// DeleteByPrefix implements [ProcessWithKey].
func (s *BaseServiceKvStore) DeleteByPrefix(prefix string) error {
	if s.kvstore == nil {
		return fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}
	return s.kvstore.DeleteByPrefix(s.s.Name(), prefix)
}

// List implements ProcessWithKey.
func (s *BaseServiceKvStore) List(opts ...ListOption) ([]KeyValue, error) {

	opt := &ListOptions{}
	for _, o := range opts {
		o(opt)
	}

	items, err := s.kvstore.List(s.s.Name(), kvstore.WithPrefix(opt.Prefix))
	if err != nil {
		return nil, err
	}

	kvs := []KeyValue{}
	for _, item := range items {
		kvs = append(kvs, KeyValue(item))
	}

	return kvs, nil

}

// ProcessWithKey implements Service.
func (s *BaseServiceKvStore) ProcessWithKey(ctx context.Context, key string, fn func(ctx context.Context) error) error {

	kv, err := s.Get(key)
	if err == nil && kv.Value != nil {
		slog.Debug("service process with key ,alreay handled", "service", s.s.Name(), "key", key)
		return nil
	}

	err = fn(ctx)
	if err != nil {
		slog.Error("service process with key ", "service", s.s.Name(), "key", key, "err", err)
		return err
	}
	if err = s.Put(key, time.Now().String()); err != nil {
		slog.Error("service process with key , put err", "service", s.s.Name(), "key", key, "err", err)
		return err
	}

	return nil

}

// Delete implements Service.
func (s *BaseServiceKvStore) Delete(key string) error {
	if s.kvstore == nil {
		return fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}
	return s.kvstore.Delete(s.s.Name(), key)
}

// Get implements Service.
func (s *BaseServiceKvStore) Get(key string) (KeyValue, error) {
	if s.kvstore == nil {
		return KeyValue{}, fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}

	kv, err := s.kvstore.Get(s.s.Name(), key)
	if err != nil {
		return KeyValue{}, err
	}
	return KeyValue(kv), nil
}

// Put implements Service.
func (s *BaseServiceKvStore) Put(key string, value any) error {
	if s.kvstore == nil {
		return fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}

	bb, _ := json.Marshal(value)

	return s.kvstore.Put(s.s.Name(), key, bb)
}

func NewBaseServiceKvStore(s Service, kvstore kvstore.KvStore) ProcessWithKey {
	return &BaseServiceKvStore{
		kvstore: kvstore,
		s:       s,
	}
}
