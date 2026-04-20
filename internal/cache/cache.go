package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, data []byte, ttl time.Duration)
	Get(key string) ([]byte, bool)
}

type CacheItem struct {
	Data      []byte
	ExpiresAt time.Time
}

type memoryCache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewCache() Cache {
	return &memoryCache{
		items: make(map[string]CacheItem),
	}
}

func (c *memoryCache) Set(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *memoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}
