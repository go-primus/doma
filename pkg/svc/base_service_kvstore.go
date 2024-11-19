package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/primus/primus/pkg/kvstore"
	"github.com/primus/primus/pkg/logger"
)

type BaseServiceKvStore struct {
	kvstore kvstore.KvStore
	s       Service
}

// ProcessWithKey implements Service.
func (s *BaseServiceKvStore) ProcessWithKey(ctx context.Context, key string, fn func(ctx context.Context) error) error {

	value, err := s.Get([]byte(key))
	if err == nil && value != nil {
		logger.L().Debugf("service [%s] process with key %s ,alreay handled", s.s.Name(), key)
		return nil
	}

	err = fn(ctx)
	if err != nil {
		logger.L().Errorf("service [%s] process with key %s, err:%v", s.s.Name(), key, err)
		return err
	}
	if err = s.Put([]byte(key), []byte(time.Now().String())); err != nil {
		logger.L().Errorf("service [%s] process with key %s, put err: %v", s.s.Name(), key, err)
		return err
	}

	return nil

}

// Delete implements Service.
func (s *BaseServiceKvStore) Delete(key []byte) error {
	if s.kvstore == nil {
		return fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}
	return s.kvstore.Delete([]byte(s.s.Name()), key)
}

// Get implements Service.
func (s *BaseServiceKvStore) Get(key []byte) ([]byte, error) {
	if s.kvstore == nil {
		return []byte{}, fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}

	return s.kvstore.Get([]byte(s.s.Name()), key)
}

// Put implements Service.
func (s *BaseServiceKvStore) Put(key []byte, value []byte) error {
	if s.kvstore == nil {
		return fmt.Errorf("service [%s] kv store is nil ", s.s.Name())
	}

	return s.kvstore.Put([]byte(s.s.Name()), key, value)
}

func NewBaseServiceKvStore(s Service, kvstore kvstore.KvStore) ProcessWithKey {
	return &BaseServiceKvStore{
		kvstore: kvstore,
		s:       s,
	}
}
