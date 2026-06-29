package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try to get request ID from client header (optional)
		requestID := c.GetHeader("X-Request-ID")

		// If not provided, generate a new one
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Save in Gin context
		c.Set(RequestIDKey, requestID)

		// Return in response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		// 5. continua a cadeia
		c.Next()
	}
}