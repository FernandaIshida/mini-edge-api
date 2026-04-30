package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
	val, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// key does not exist
			return nil, 0, false
		}

		log.Println("[REDIS ERROR - GET]:", err)
		return nil, 0, false
	}

	var item CacheItem
	if err := json.Unmarshal(val, &item); err != nil {
		log.Println("[REDIS ERROR - UNMARSHAL]:", err)
		return nil, 0, false
	}

	log.Println("[REDIS HIT]:", key)
	return item.Data, item.StatusCode, true
}

func (r *RedisCache) Set(key string, value []byte, status int, ttl time.Duration) {
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

	log.Println("[REDIS SET]:", key)
}
