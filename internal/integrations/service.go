package integrations

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/FernandaIshida/mini-edge-api/internal/metrics"
)

type Service struct {
	cache      cache.Cache
	redisCache *cache.RedisCache
	client     *http.Client
}

func NewService(
	cache cache.Cache,
	redisCache *cache.RedisCache,
	client *http.Client,
) *Service {
	return &Service{
		cache:      cache,
		redisCache: redisCache,
		client:     client,
	}
}

func (s *Service) Check(requestID, name string) ([]byte, int, string, error) {
	key := "integration:" + name

	//L1 (memory) cache check
	if data, status, found := s.cache.Get(key); found {
		log.Printf("[request_id=%s] CACHE HIT L1: %s", requestID, key)

		return data, status, "HIT", nil
	}
	//L2 redis
	if data, status, found := s.redisCache.Get(key); found {
		log.Printf("[request_id=%s] CACHE HIT L2: %s", requestID, key)

		//Warm up L1 cache after Redis hit
		s.cache.Set(key, data, status, 30*time.Second)

		return data, status, "HIT-REDIS", nil
	}

	log.Printf("[request_id=%s] CACHE MISS: %s", requestID, key)

	integration, exists := Registry[name]
	if !exists {
		log.Printf("[request_id=%s] INTEGRATION NOT FOUND: %s", requestID, name)

		return nil, 0, "", errors.New("integration not found")
	}

	url := integration.URL

	// Create request to external integration
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, "", err
	}

	start := time.Now()

	resp, err := s.client.Do(req)

	metrics.ExternalAPIDuration.WithLabelValues(name).Observe(time.Since(start).Seconds())

	log.Printf("[request_id=%s] EXTERNAL CALL: %s (%s)", requestID, name, url)

	if err != nil {
		return nil, 0, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}

	if resp.StatusCode == http.StatusOK {

		s.cache.Set(
			key,
			body,
			resp.StatusCode,
			30*time.Second,
		)

		s.redisCache.Set(
			key,
			body,
			resp.StatusCode,
			30*time.Second,
		)
	}

	return body, resp.StatusCode, "MISS", nil
}
