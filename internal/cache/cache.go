package cache

import (
	"sync"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/metrics"
)

type Cache interface {
	Set(key string, data []byte, status int, ttl time.Duration)
	Get(key string) ([]byte, int, bool)
}

type CacheItem struct {
	Data       []byte    `json:"data"`
	StatusCode int       `json:"status_code"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type memoryCache struct {
	mu       sync.RWMutex
	items    map[string]CacheItem
	interval time.Duration
	stop     chan struct{}
}

func NewCache(cleanupInterval time.Duration) Cache {
	c := &memoryCache{
		items:    make(map[string]CacheItem),
		interval: cleanupInterval,
		stop:     make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go c.startCleanup()
	}

	return c
}

func (c *memoryCache) Set(key string, data []byte, status int, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Data:       data,
		StatusCode: status,
		ExpiresAt:  time.Now().Add(ttl),
	}

	metrics.CacheWrites.Inc()

}

func (c *memoryCache) Get(key string) ([]byte, int, bool) {
	start := time.Now()

	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		metrics.CacheMiss.Inc()
		return nil, 0, false
	}

	now := time.Now()
	//Lazy eviction: remove expired items on access
	if now.After(item.ExpiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()

		metrics.CacheMiss.Inc()

		metrics.CacheDuration.WithLabelValues("miss").Observe(time.Since(start).Seconds())

		return nil, 0, false
	}

	metrics.CacheHits.Inc()

	metrics.CacheDuration.WithLabelValues("hit").Observe(time.Since(start).Seconds())

	return item.Data, item.StatusCode, true
}

func (c *memoryCache) startCleanup() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()

			c.mu.Lock()
			for key, item := range c.items {
				if now.After(item.ExpiresAt) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()

		case <-c.stop:
			return
		}
	}
}
