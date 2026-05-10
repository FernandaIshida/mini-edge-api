package main

import (
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/api"
	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/FernandaIshida/mini-edge-api/internal/metrics"
	"github.com/FernandaIshida/mini-edge-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	redisCache := cache.NewRedisCache("localhost:6379")

	cacheService := cache.NewCache(1 * time.Second)

	metrics.Init()

	r := gin.Default()

	r.Use(middleware.RateLimiter())
	r.Use(middleware.MetricsMiddleware())

	api.SetupRoutes(r, cacheService, redisCache, httpClient)

	r.Run(":8080")
}
