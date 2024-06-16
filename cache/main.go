package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Cache struct for storing AI predictions
type Cache struct {
	mu     sync.RWMutex
	result map[string]int
}

// NewCache creates a new Cache instance and starts the invalidation goroutine
func NewCache(ctx context.Context) *Cache {
	result := make(map[string]int)
	cache := &Cache{
		result: result,
	}
	// Generate initial ai_predict value
	cache.result["ai_predict"] = aiPredict()
	go cache.invalidate(ctx)
	return cache
}

// invalidate clears the cache every 2 seconds
func (c *Cache) invalidate(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			fmt.Println("Invalidated")
			c.mu.Lock()
			for key := range c.result {
				delete(c.result, key)
			}
			c.mu.Unlock()
		}
	}
}

// GetValue returns the cached value or generates a new one if not present
func (c *Cache) GetValue(key string) int {
	c.mu.RLock()
	value, ok := c.result[key]
	c.mu.RUnlock()

	if ok {
		fmt.Println("CACHE hit")
		return value
	} else {
		fmt.Println("CACHE miss")
		newAiPredict := aiPredict()
		c.mu.Lock()
		c.result[key] = newAiPredict
		c.mu.Unlock()
		return newAiPredict
	}
}

// aiPredict simulates an AI prediction
func aiPredict() int {
	time.Sleep(2 * time.Second)
	return rand.Intn(100)
}

func main() {
	// Context for cache invalidation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := NewCache(ctx)

	// HTTP handler for /predict endpoint
	http.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		result := cache.GetValue("ai_predict")
		fmt.Fprintf(w, "{\"result\": %d}", result)
	})

	// Start HTTP server
	fmt.Println("Server starting on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
