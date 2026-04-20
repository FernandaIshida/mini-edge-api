package handlers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

func ExternalHandler(cache cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := ctx.Request.Context()
		key := "external:posts"

		ctx.Header("Cache-Control", "public, max-age=30")

		if data, found := cache.Get(key); found {
			log.Println("CACHE HIT")
			ctx.Header("X-Cache", "HIT")
			ctx.Data(http.StatusOK, "application/json", data)
			return
		}

		ctx.Header("X-Cache", "MISS")

		url := "https://jsonplaceholder.typicode.com/posts"

		// Create a new request with the context
		req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to create request"})
			return
		}

		// Use the default HTTP client to make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "external request failed"})
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "failed to read response"})
			return
		}

		cache.Set(key, body, 30*time.Second)

		// Return the external API response directly to the client
		ctx.Data(http.StatusOK, "application/json", body)
	}
}
