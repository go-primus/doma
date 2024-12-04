package kvstore

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"
)

type KvStore interface {
	Create(ctx context.Context, key string, value []byte) (int64, error)
	Delete(ctx context.Context, key string) error
	Update(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) (any, error)
	List(ctx context.Context, prefix, startkey string, limit int64) ([]any, error)
	Count(ctx context.Context, prefix, startkey string) (int64, error)

	//
	Watch()
}

type Watcher interface {
	Watch(ctx context.Context, key string) WatchResult
}

type Event struct {
}

type WatchResult struct {
	Events <-chan []*Event
}

type kvStore struct {
	db *bbolt.DB
}

// Delete implements ConfigStore.
func (s *kvStore) Delete(bucket []byte, key []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		return bucket.Delete(key)

	})
}

// Get implements ConfigStore.
func (s *kvStore) Get(bucket []byte, key []byte) ([]byte, error) {
	value := []byte{}
	err := s.db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("bucket is null")
		}
		value = bucket.Get(key)
		return nil
	})
	return value, err
}

// Put implements ConfigStore.
func (s *kvStore) Put(bucket []byte, key, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		return bucket.Put(key, value)
	})
}

func NewKvStore(db *bbolt.DB) KvStore {
	return &kvStore{
		db: db,
	}
}
