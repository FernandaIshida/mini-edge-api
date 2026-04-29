package main

import (
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/api"
	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/FernandaIshida/mini-edge-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.RateLimiter())

	cacheService := cache.NewCache(1 * time.Second)

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	api.SetupRoutes(r, cacheService, httpClient)

	r.Run(":8080")
}
