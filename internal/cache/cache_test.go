package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(0)

	cache.Set("key", []byte("value"), time.Second)

	data, ok := cache.Get("key")

	if !ok {
		t.Fatal("Expected key to exist")
	}

	if string(data) != "value" {
		t.Fatalf("Expected value to be 'value', got '%s'", data)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(0)

	cache.Set("key", []byte("value"), 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, ok := cache.Get("key")

	if ok {
		t.Fatal("Expected key to be expired")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	cache.Set("key|", []byte("value"), 300*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	_, ok := cache.Get("key")

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
			cache.Set(key, []byte(fmt.Sprintf("value%d", i)), time.Second)
			cache.Get(key)
		}(i)
	}

	wg.Wait()
}
