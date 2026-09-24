package store

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	data     map[string]any
	watchers map[string]map[chan WatchEvent]struct{}
}

var _ Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:     make(map[string]any),
		watchers: make(map[string]map[chan WatchEvent]struct{}),
	}
}

func (m *MemoryStore) Put(ctx context.Context, key string, value any) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	exists := false
	if _, ok := m.data[key]; ok {
		exists = true
	}

	m.data[key] = value

	eventType := EventCreate
	if exists {
		eventType = EventUpdate
	}
	m.notifyWatchers(key, eventType)

	return nil
}

func (m *MemoryStore) Get(ctx context.Context, key string) (KeyValue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	if !ok {
		return KeyValue{}, ErrKeyNotFound
	}

	return KeyValue{Key: key, Value: value}, nil
}

func (m *MemoryStore) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return ErrKeyNotFound
	}

	delete(m.data, key)
	m.notifyWatchers(key, EventDelete)

	return nil
}

func (m *MemoryStore) DeletePrefix(ctx context.Context, prefix string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if prefix == "" {
		return nil
	}

	for key := range m.data {
		if matchPrefix(key, prefix) {
			delete(m.data, key)
			m.notifyWatchers(key, EventDelete)
		}
	}

	return nil
}

func (m *MemoryStore) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[key]
	return ok, nil
}

func (m *MemoryStore) Update(ctx context.Context, key string, value any) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return ErrKeyNotFound
	}

	m.data[key] = value
	m.notifyWatchers(key, EventUpdate)

	return nil
}

func (m *MemoryStore) List(ctx context.Context, opts ListOptions) (ListResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limit := opts.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	var keys []string
	for k := range m.data {
		if opts.Prefix != "" && !matchPrefix(k, opts.Prefix) {
			continue
		}
		keys = append(keys, k)
	}

	startIdx := 0
	if opts.Cursor != "" {
		for i, k := range keys {
			if k == opts.Cursor {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(keys) {
		endIdx = len(keys)
	}

	if startIdx >= len(keys) {
		return ListResult{Items: []KeyValue{}}, nil
	}

	items := make([]KeyValue, 0, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		items = append(items, KeyValue{
			Key:   keys[i],
			Value: m.data[keys[i]],
		})
	}

	cursor := ""
	if endIdx < len(keys) {
		cursor = keys[endIdx-1]
	}

	return ListResult{Items: items, Cursor: cursor}, nil
}

func (m *MemoryStore) Watch(ctx context.Context, prefix string) (<-chan WatchEvent, error) {
	ch := make(chan WatchEvent, 100)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.watchers[prefix] == nil {
		m.watchers[prefix] = make(map[chan WatchEvent]struct{})
	}
	m.watchers[prefix][ch] = struct{}{}

	go m.cleanupWatcher(ctx, prefix, ch)

	return ch, nil
}

func (m *MemoryStore) notifyWatchers(key string, eventType EventType) {
	event := WatchEvent{
		Type: eventType,
		KeyValue: KeyValue{
			Key:   key,
			Value: m.data[key],
		},
	}

	if eventType == EventDelete {
		event.Value = nil
	}

	for prefix, watchers := range m.watchers {
		if matchPrefix(key, prefix) {
			for ch := range watchers {
				select {
				case ch <- event:
				default:
				}
			}
		}
	}
}

func (m *MemoryStore) cleanupWatcher(ctx context.Context, prefix string, ch chan WatchEvent) {
	select {
	case <-ctx.Done():
		slog.Warn("context done")
		// case <-ch:
		// 	slog.Warn("event ch receive")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if watchers, ok := m.watchers[prefix]; ok {
		delete(watchers, ch)
		close(ch)
		if len(watchers) == 0 {
			delete(m.watchers, prefix)
		}
	}
}

// matchPrefix 检查给定的 key 是否与指定的模式匹配。
// 模式支持通配符 '*'，它可以匹配除 ':' 以外的任意字符序列。
// 参数：
//   - key: 要检查是否匹配模式的字符串。
//   - pattern: 要匹配的模式，其中 '*' 表示通配符。
//
// 返回值：
//   - 如果 key 与模式匹配则返回 true，否则返回 false。
func matchPrefix(key, pattern string) bool {
	// 如果模式为空，则认为匹配任何 key。
	if pattern == "" {
		return true
	}

	// 如果模式不包含通配符，则检查 key 是否以该模式开头。
	if !strings.Contains(pattern, "*") {
		return strings.HasPrefix(key, pattern)
	}

	// 将模式转换为正则表达式：
	// 1. 转义点号 '.'，使其作为字面量字符处理。
	// 2. 将 '*' 替换为 "[^:]*"，表示匹配除 ':' 外的任意字符序列。
	// 3. 添加锚点 "^" 和 "$"，确保从头到尾完全匹配。
	re := strings.ReplaceAll(pattern, ".", "\\.")
	re = strings.ReplaceAll(re, "*", "[^:]*")
	re = "^" + re + "$"

	// 执行正则匹配并返回结果。
	matched, _ := regexp.MatchString(re, key)
	return matched
}
