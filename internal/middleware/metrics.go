package middleware

import (
	"strconv"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/metrics"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		metrics.HttpRequestTotal.WithLabelValues(
			c.FullPath(),
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			c.FullPath(),
		).Observe(duration)
	}
}
