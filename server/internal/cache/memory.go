package cache

import (
	"context"
	"sync"
	"time"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

type memoryItem struct {
	value     []byte
	expiresAt time.Time
}

func NewMemory() *Memory {
	return &Memory{items: map[string]memoryItem{}}
}

func (m *Memory) Ping(_ context.Context) error {
	return nil
}

func (m *Memory) Get(_ context.Context, key string) ([]byte, bool, error) {
	m.mu.RLock()
	item, ok := m.items[key]
	m.mu.RUnlock()
	if !ok {
		return nil, false, nil
	}
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		_ = m.Delete(context.Background(), key)
		return nil, false, nil
	}
	return append([]byte(nil), item.value...), true, nil
}

func (m *Memory) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	item := memoryItem{value: append([]byte(nil), value...)}
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	m.mu.Lock()
	m.items[key] = item
	m.mu.Unlock()
	return nil
}

func (m *Memory) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
	return nil
}

func (m *Memory) Close() error {
	return nil
}
