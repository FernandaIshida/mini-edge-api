package cache

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(0)

	cache.Set("key", []byte("value"), http.StatusOK, time.Second)

	data, status, ok := cache.Get("key")

	if !ok {
		t.Fatal("Expected key to exist")
	}

	if status != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", status)
	}

	if string(data) != "value" {
		t.Fatalf("Expected value to be 'value', got '%s'", data)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(0)

	cache.Set("key", []byte("value"), http.StatusOK, 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, _, ok := cache.Get("key")

	if ok {
		t.Fatal("Expected key to be expired")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	cache.Set("key|", []byte("value"), http.StatusOK, 300*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	_, _, ok := cache.Get("key")

	if ok {
		t.Fatal("Expected key to be removed by cleanup")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(0)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("key-%d", i)

			cache.Set(
				key,
				[]byte(fmt.Sprintf("value%d", i)), http.StatusOK,
				time.Second,
			)

			cache.Get(key)
		}(i)
	}

	wg.Wait()
}

func TestCacheClose(t *testing.T) {
	cache := NewCache(time.Second)
	cache.Close()
}

func TestCloseIsIdempotent(t *testing.T) {
	cache := NewCache(time.Second)

	cache.Close()
	cache.Close() // Should not panic or cause any issues
}

func TestCacheWithoutCleanupCanClose(t *testing.T) {
	cache := NewCache(0)
	cache.Close()
}
