package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	requests int
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
	limit   = 5           // 5 requests per minute
	window  = time.Minute // 60 seconds
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		cl, exists := clients[ip]
		if !exists {
			clients[ip] = &client{
				requests: 1,
				lastSeen: time.Now(),
			}
			mu.Unlock()
			c.Next()
			return
		}

		// Reset the count if the time window has passed
		if time.Since(cl.lastSeen) > window {
			cl.requests = 0
			cl.lastSeen = time.Now()
		}

		cl.requests++

		if cl.requests > limit {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}
