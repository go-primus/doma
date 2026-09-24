package kvstore_test

import (
	"errors"
	"testing"

	"github.com/go-primus/doma/pkg/kvstore"
)

func newTestStore(t *testing.T) kvstore.KvStore {
	store, err := kvstore.Open("bbolt", t.TempDir()+"/test.db")
	if err != nil {
		t.Fatalf("failed to open kvstore: %v", err)
	}
	return store
}

func TestKvStore_Put(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	err := store.Put("bucket", "key1", []byte("value1"))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	err = store.Put("bucket", "key2", []byte("value2"))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
}

func TestKvStore_Get(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key", []byte("value"))

	kv, err := store.Get("bucket", "key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if kv.Key != "key" {
		t.Errorf("expected key 'key', got '%s'", kv.Key)
	}
	if string(kv.Value) != "value" {
		t.Errorf("expected value 'value', got '%s'", string(kv.Value))
	}
}

func TestKvStore_Get_KeyNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "somekey", []byte("somevalue"))

	_, err := store.Get("bucket", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key, got nil")
	}
	if !errors.Is(err, kvstore.ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestKvStore_Get_BucketNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	_, err := store.Get("nonexistent", "key")
	if err == nil {
		t.Error("expected error for nonexistent bucket, got nil")
	}
	if !errors.Is(err, kvstore.ErrBucketNotFound) {
		t.Errorf("expected ErrBucketNotFound, got %v", err)
	}
}

func TestKvStore_Delete(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key", []byte("value"))

	err := store.Delete("bucket", "key")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get("bucket", "key")
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestKvStore_Delete_KeyNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "somekey", []byte("somevalue"))

	err := store.Delete("bucket", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key, got nil")
	}
	if !errors.Is(err, kvstore.ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestKvStore_Delete_BucketNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	err := store.Delete("nonexistent", "key")
	if err == nil {
		t.Error("expected error for nonexistent bucket, got nil")
	}
	if !errors.Is(err, kvstore.ErrBucketNotFound) {
		t.Errorf("expected ErrBucketNotFound, got %v", err)
	}
}

func TestKvStore_DeleteByPrefix(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "app:user:1", []byte("data1"))
	store.Put("bucket", "app:user:2", []byte("data2"))
	store.Put("bucket", "app:profile:1", []byte("data3"))
	store.Put("bucket", "other:data", []byte("data4"))

	err := store.DeleteByPrefix("bucket", "app:user:")
	if err != nil {
		t.Fatalf("DeleteByPrefix failed: %v", err)
	}

	kv, _ := store.Get("bucket", "app:user:1")
	if kv.Key != "" {
		t.Error("expected app:user:1 to be deleted")
	}

	kv, _ = store.Get("bucket", "app:user:2")
	if kv.Key != "" {
		t.Error("expected app:user:2 to be deleted")
	}

	kv, _ = store.Get("bucket", "app:profile:1")
	if kv.Key == "" {
		t.Error("expected app:profile:1 to remain (different prefix)")
	}

	kv, _ = store.Get("bucket", "other:data")
	if kv.Key == "" {
		t.Error("expected other:data to remain")
	}
}

func TestKvStore_DeleteByPrefix_NoMatch(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key1", []byte("value1"))

	err := store.DeleteByPrefix("bucket", "nonexistent")
	if err != nil {
		t.Fatalf("DeleteByPrefix failed: %v", err)
	}

	kv, _ := store.Get("bucket", "key1")
	if kv.Key != "key1" {
		t.Error("expected key1 to remain")
	}
}

func TestKvStore_DeleteByPrefix_BucketNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	err := store.DeleteByPrefix("nonexistent", "prefix")
	if err == nil {
		t.Error("expected error for nonexistent bucket, got nil")
	}
	if !errors.Is(err, kvstore.ErrBucketNotFound) {
		t.Errorf("expected ErrBucketNotFound, got %v", err)
	}
}

func TestKvStore_List(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "a:1", []byte("1"))
	store.Put("bucket", "a:2", []byte("2"))
	store.Put("bucket", "b:1", []byte("3"))

	items, err := store.List("bucket")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestKvStore_List_WithPrefix(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "app:user:1", []byte("1"))
	store.Put("bucket", "app:user:2", []byte("2"))
	store.Put("bucket", "app:profile:1", []byte("3"))
	store.Put("bucket", "other:data", []byte("4"))

	items, err := store.List("bucket", kvstore.WithPrefix("app:user:"))
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	for _, item := range items {
		if len(item.Key) < 10 {
			t.Errorf("key %s should match prefix app:user:", item.Key)
		}
	}
}

func TestKvStore_List_WithLimit(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	for i := 1; i <= 10; i++ {
		key := "key" + string(rune('0'+i))
		store.Put("bucket", key, []byte("value"))
	}
	store.Put("bucket", "key1", []byte("value1"))
	store.Put("bucket", "key2", []byte("value2"))
	store.Put("bucket", "key3", []byte("value3"))

	items, err := store.List("bucket", kvstore.WithLimit(2))
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items with limit, got %d", len(items))
	}
}

func TestKvStore_List_WithOffset(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key1", []byte("1"))
	store.Put("bucket", "key2", []byte("2"))
	store.Put("bucket", "key3", []byte("3"))

	items, err := store.List("bucket", kvstore.WithOffset(1))
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items with offset 1, got %d", len(items))
	}
}

func TestKvStore_List_WithReverse(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key1", []byte("1"))
	store.Put("bucket", "key2", []byte("2"))
	store.Put("bucket", "key3", []byte("3"))

	items, err := store.List("bucket", kvstore.WithReverse())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
	if items[0].Key != "key3" || items[2].Key != "key1" {
		t.Error("expected reverse order: key3, key2, key1")
	}
}

func TestKvStore_List_BucketNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	_, err := store.List("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent bucket, got nil")
	}
	if !errors.Is(err, kvstore.ErrBucketNotFound) {
		t.Errorf("expected ErrBucketNotFound, got %v", err)
	}
}

func TestKvStore_List_CombinedOptions(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	for i := 1; i <= 10; i++ {
		store.Put("bucket", "app:user:"+string(rune(i)), []byte("value"))
	}
	store.Put("bucket", "app:admin:1", []byte("admin"))
	store.Put("bucket", "other:data", []byte("other"))

	items, err := store.List("bucket",
		kvstore.WithPrefix("app:user:"),
		kvstore.WithOffset(2),
		kvstore.WithLimit(3),
	)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestKvStore_Close(t *testing.T) {
	store := newTestStore(t)

	err := store.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	err = store.Put("bucket", "key", []byte("value"))
	if err == nil {
		t.Error("expected error after close, got nil")
	}
}

func TestKvStore_AutoCreateBucket(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	err := store.Put("newbucket", "key", []byte("value"))
	if err != nil {
		t.Fatalf("Put to new bucket failed: %v", err)
	}

	kv, err := store.Get("newbucket", "key")
	if err != nil {
		t.Fatalf("Get from new bucket failed: %v", err)
	}
	if string(kv.Value) != "value" {
		t.Errorf("expected 'value', got '%s'", string(kv.Value))
	}
}

func TestKvStore_EmptyBucket(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	store.Put("bucket", "key", []byte("value"))

	store.Delete("bucket", "key")

	items, err := store.List("bucket")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestKvStore_DriverRegistry(t *testing.T) {
	driver, ok := kvstore.Get("bbolt")
	if !ok {
		t.Error("expected bbolt driver to be registered")
	}
	if driver.Name() != "bbolt" {
		t.Errorf("expected driver name 'bbolt', got '%s'", driver.Name())
	}

	names := kvstore.Names()
	if len(names) != 1 || names[0] != "bbolt" {
		t.Errorf("expected ['bbolt'], got %v", names)
	}
}
