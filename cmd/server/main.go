package main

import (
	"github.com/FernandaIshida/mini-edge-api/internal/api"
	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/FernandaIshida/mini-edge-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.RateLimiter())

	cache := cache.NewCache()

	api.SetupRoutes(r, cache)

	r.Run(":8080")
}
