package api

import (
	"net/http"

	"github.com/FernandaIshida/mini-edge-api/internal/api/handlers"
	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/FernandaIshida/mini-edge-api/internal/integrations"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(r *gin.Engine, cache cache.Cache, redisCache *cache.RedisCache, httpClient *http.Client) {
	r.GET("/health", handlers.HealthHandler)

	integrationService := integrations.NewService(
		cache,
		redisCache,
		httpClient,
	)

	integrationHandler := handlers.IntegrationHandler(
		integrationService,
	)
	r.GET("/integrations/:name", integrationHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
