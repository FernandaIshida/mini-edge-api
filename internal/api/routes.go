package api

import (
	"github.com/FernandaIshida/mini-edge-api/internal/api/handlers"
	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cache cache.Cache) {
	r.GET("health", handlers.HealthHandler)
	r.GET("/products", handlers.ProductsHandler(cache))
	r.GET("/external-data", handlers.ExternalHandler(cache))
}
