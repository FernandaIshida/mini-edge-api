package main

import (
	"github.com/FernandaIshida/mini-edge-api/internal/api"
	"github.com/FernandaIshida/mini-edge-api/internal/limiter"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(limiter.RateLimiter())

	api.SetupRoutes(r)

	r.Run(":8080")
}
