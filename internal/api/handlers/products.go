package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

type Products struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var mockProducts = []Products{
	{ID: 1, Name: "Product A", Price: 10.99},
	{ID: 2, Name: "Product B", Price: 19.99},
	{ID: 3, Name: "Product C", Price: 5.49},
}

func ProductsHandler(cache cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "products:list"

		if data, found := cache.Get(key); found {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", data)
			return
		}

		c.Header("X-Cache", "MISS")

		time.Sleep(200 * time.Millisecond)

		jsonData, err := json.Marshal(mockProducts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode"})
			return
		}

		cache.Set(key, jsonData, 30*time.Second)

		c.Data(http.StatusOK, "application/json", jsonData)
	}
}
