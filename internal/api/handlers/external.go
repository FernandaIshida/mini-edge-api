package handlers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

func ExternalHandler(cache cache.Cache, redisCache *cache.RedisCache, client *http.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := ctx.Request.Context()
		key := "external:" + ctx.Request.URL.String()

		//L1 (memory) cache check
		if data, status, found := cache.Get(key); found {
			log.Printf("CACHE HIT: %s", key)
			ctx.Header("X-Cache", "HIT")
			ctx.Data(status, "application/json", data)
			return
		}
		//L2 redis
		if data, status, found := redisCache.Get(key); found {
			log.Printf("CACHE HIT (Redis): %s", key)

			redisCache.Set(key, data, status, 30*time.Second)
			ctx.Header("X-Cache", "HIT-REDIS")
			ctx.Data(status, "application/json", data)
			return
		}

		log.Printf("CACHE MISS: %s", key)
		ctx.Header("X-Cache", "MISS")

		url := "https://jsonplaceholder.typicode.com/posts"

		// Create a new request with the context
		req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "external request failed"})
			return
		}

		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
			return
		}

		if resp.StatusCode == http.StatusOK {
			ctx.Header("Cache-Control", "public, max-age=30")

			cache.Set(key, body, resp.StatusCode, 30*time.Second)
			redisCache.Set(key, body, resp.StatusCode, 30*time.Second)
		} else {
			ctx.Header("Cache-Control", "no-store")
		}

		// Return the external API response directly to the client
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}
		ctx.Data(resp.StatusCode, contentType, body)
	}
}
