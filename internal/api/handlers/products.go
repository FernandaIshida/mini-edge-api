package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/FernandaIshida/mini-edge-api/internal/cache"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var mockProducts = []Product{
	{ID: 1, Name: "Product A", Price: 10.99},
	{ID: 2, Name: "Product B", Price: 19.99},
	{ID: 3, Name: "Product C", Price: 5.49},
}

func ProductsHandler(cache cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "products:list"
		c.Header("Cache-Control", "public, max-age=30")

		if data, status, found := cache.Get(key); found {
			c.Header("X-Cache", "HIT")
			log.Printf("CACHE HIT: %s", key)
			c.Data(status, "application/json", data)
			return
		}

		log.Printf("CACHE MISS: %s", key)
		c.Header("X-Cache", "MISS")

		time.Sleep(200 * time.Millisecond)

		jsonData, err := json.Marshal(mockProducts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode"})
			return
		}

		cache.Set(key, jsonData, http.StatusOK, 30*time.Second)

		c.Data(http.StatusOK, "application/json", jsonData)
	}
}
