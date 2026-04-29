package cache

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func BenchmarkCache_GetHit(b *testing.B) {
	c := NewCache(0)

	key := "key"
	value := []byte("value")

	c.Set(key, value, http.StatusOK, time.Minute)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = c.Get(key)
	}
}

func BenchmarkCache_GetMiss(b *testing.B) {
	c := NewCache(0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = c.Get("missing-key")
	}
}

func BenchmarkCache_Set(b *testing.B) {
	c := NewCache(0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte("value")
		c.Set(key, value, http.StatusOK, time.Minute)
	}
}
