package handlers

import (
	"net/http"

	"github.com/FernandaIshida/mini-edge-api/internal/integrations"
	"github.com/gin-gonic/gin"
)

func IntegrationHandler(service *integrations.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")

		requestID := ctx.GetString("request_id")

		data, status, cacheStatus, err := service.Check(requestID, name)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Header("X-Cache", cacheStatus)
		ctx.Header("X-Request-ID", requestID)
		ctx.Data(status, "application/json", data)
	}
}
