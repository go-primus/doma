package kvstore

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

func init() {
	Register(&bboltDriver{})
}

type bboltDriver struct{}

func (d *bboltDriver) Name() string {
	return "bbolt"
}

func (d *bboltDriver) Open(path string) (KvStore, error) {

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}
	slog.Debug("bbolt driver open ", "path", path)
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &bboltStore{db: db}, nil
}

type bboltStore struct {
	db *bbolt.DB
}

func (s *bboltStore) Close() error {
	return s.db.Close()
}

func (s *bboltStore) Put(bucket, key string, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}

func (s *bboltStore) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		v := b.Get([]byte(key))
		if v == nil {
			return ErrKeyNotFound
		}
		return b.Delete([]byte(key))
	})
}

func (s *bboltStore) DeleteByPrefix(bucket, prefix string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		prefixBytes := []byte(prefix)

		for k, _ := c.First(); k != nil; {
			if bytes.HasPrefix(k, prefixBytes) {
				if err := c.Delete(); err != nil {
					return err
				}
				k, _ = c.Next()
			} else {
				k, _ = c.Next()
			}
		}
		return nil
	})
}

func (s *bboltStore) Get(bucket, key string) (KeyValue, error) {
	var kv KeyValue
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrKeyNotFound
		}
		v := b.Get([]byte(key))
		if v == nil {
			return ErrKeyNotFound
		}
		kv = KeyValue{Key: key, Value: v}
		return nil
	})
	return kv, err
}

func (s *bboltStore) Count(bucket string, prefix string) (int64, error) {
	var count int64
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		prefixBytes := []byte(prefix)

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if bytes.HasPrefix(k, prefixBytes) {
				count++
			}
		}
		return nil
	})
	return count, err
}

func (s *bboltStore) List(bucket string, opts ...ListOption) ([]KeyValue, error) {
	opt := &ListOptions{}
	for _, o := range opts {
		o(opt)
	}

	var items []KeyValue
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}

		c := b.Cursor()
		prefixBytes := []byte(opt.Prefix)

		var iterate func(k, v []byte) error
		if opt.Reverse {
			iterate = func(k, v []byte) error {
				if !bytes.HasPrefix(k, prefixBytes) {
					return nil
				}
				items = append(items, KeyValue{Key: string(k), Value: v})
				return nil
			}
		} else {
			iterate = func(k, v []byte) error {
				if !bytes.HasPrefix(k, prefixBytes) {
					return nil
				}
				items = append(items, KeyValue{Key: string(k), Value: v})
				return nil
			}
		}

		if opt.Reverse {
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				if err := iterate(k, v); err != nil {
					return err
				}
				if opt.Limit > 0 && len(items) >= opt.Limit+opt.Offset {
					break
				}
			}
		} else {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if err := iterate(k, v); err != nil {
					return err
				}
				if opt.Limit > 0 && len(items) >= opt.Limit+opt.Offset {
					break
				}
			}
		}

		if opt.Offset > 0 && len(items) > opt.Offset {
			items = items[opt.Offset:]
		}
		if opt.Limit > 0 && len(items) > opt.Limit {
			items = items[:opt.Limit]
		}

		return nil
	})
	return items, err
}
