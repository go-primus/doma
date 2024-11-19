package kvstore

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type KvStore interface {
	Put(bucket []byte, key, value []byte) error
	Delete(bucket []byte, key []byte) error
	Get(bucket []byte, key []byte) ([]byte, error)
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
