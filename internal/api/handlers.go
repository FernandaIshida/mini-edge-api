package api

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

var c = cache.NewCache()

func DataHandler(ctx *gin.Context) {
	typeParam := ctx.Query("type")

	if typeParam == "" {
		typeParam = "default"
	}

	key := "data:type" + typeParam

	// Try to get data from cache
	if data, found := c.Get(key); found {
		ctx.Header("X-Cache", "HIT")
		ctx.Data(http.StatusOK, "application/json", data)
		return
	}

	ctx.Header("X-Cache", "MISS")

	// Simulate data fetching delay
	time.Sleep(2 * time.Second)

	// Simulated data response
	response := `{
		"source": "internal",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `",
		"items": [
			{"id": 1, "name": "Item A"},
			{"id": 2, "name": "Item B"}
		]
	}`

	// Store the response in cache with a TTL of 15 seconds
	c.Set(key, []byte(response), 15*time.Second)

	ctx.Data(http.StatusOK, "application/json", []byte(response))

}

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func ExternalHandler(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	key := "external:posts"

	ctx.Header("Cache-Control", "public, max-age=30")

	if data, found := c.Get(key); found {
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

	c.Set(key, body, 30*time.Second)

	// Return the external API response directly to the client
	ctx.Data(http.StatusOK, "application/json", body)
}
