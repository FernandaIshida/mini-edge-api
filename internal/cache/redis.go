package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/metrics"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr, // ex: "localhost:6379"
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Get(key string) ([]byte, int, bool) {
	start := time.Now()

	val, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			metrics.RedisMiss.Inc()

			metrics.RedisDuration.WithLabelValues("miss").Observe(time.Since(start).Seconds())

			// key does not exist
			return nil, 0, false
		}

		log.Println("[REDIS ERROR - GET]:", err)

		metrics.RedisMiss.Inc()

		metrics.RedisDuration.WithLabelValues("miss").Observe(time.Since(start).Seconds())

		return nil, 0, false
	}

	var item CacheItem
	if err := json.Unmarshal(val, &item); err != nil {
		log.Println("[REDIS ERROR - UNMARSHAL]:", err)

		metrics.RedisMiss.Inc()

		metrics.RedisDuration.WithLabelValues("miss").Observe(time.Since(start).Seconds())

		return nil, 0, false
	}

	log.Println("[REDIS HIT]:", key)

	metrics.RedisDuration.WithLabelValues("hit").Observe(time.Since(start).Seconds())

	return item.Data, item.StatusCode, true
}

func (r *RedisCache) Set(key string, value []byte, status int, ttl time.Duration) {
	start := time.Now()

	item := CacheItem{
		Data:       value,
		StatusCode: status,
		ExpiresAt:  time.Now().Add(ttl), // opcional (Redis já controla TTL)
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		log.Println("[REDIS ERROR - MARSHAL]:", err)
		return
	}

	err = r.client.Set(r.ctx, key, jsonData, ttl).Err()
	if err != nil {
		log.Println("[REDIS ERROR - SET]:", err)
		return
	}

	metrics.RedisWrites.Inc()

	metrics.RedisDuration.WithLabelValues("write").Observe(time.Since(start).Seconds())

	log.Println("[REDIS SET]:", key)
}
