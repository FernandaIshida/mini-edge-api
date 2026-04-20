package api

import (
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
