package api

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	r.GET("health", HealthHandler)
	r.GET("/data", DataHandler) // Support query parameter "type" for different data types
}
