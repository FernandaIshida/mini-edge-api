package main

import (
	"github.com/FernandaIshida/mini-edge-api/internal/api"
	limiter "github.com/FernandaIshida/mini-edge-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(limiter.RateLimiter())

	api.SetupRoutes(r)

	r.Run(":8080")
}
