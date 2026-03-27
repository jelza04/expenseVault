package services

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// cacheEntry stores a cached response and its expiration time.
type cacheEntry struct {
	response  string
	expiresAt time.Time
}

// AICache is a simple thread-safe in-memory cache for AI responses.
// UNIT 5: Demonstrates concurrency safety with sync.RWMutex.
type AICache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	ttl     time.Duration
}

// NewAICache creates a new AICache with the given TTL.
func NewAICache(ttl time.Duration) *AICache {
	cache := &AICache{
		entries: make(map[string]cacheEntry),
		ttl:     ttl,
	}
	// Periodically clean up expired entries in a background goroutine.
	// UNIT 5: Asynchronous background worker.
	go cache.cleanUpLoop()
	return cache
}

// Get retrieves a cached response if it exists and hasn't expired.
func (c *AICache) Get(userID int64, query string) (string, bool) {
	key := c.generateKey(userID, query)
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.response, true
}

// Set stores a response in the cache with the configured TTL.
func (c *AICache) Set(userID int64, query, response string) {
	key := c.generateKey(userID, query)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		response:  response,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// generateKey creates a unique hash for a user + query combination.
func (c *AICache) generateKey(userID int64, query string) string {
	input := fmt.Sprintf("%d:%s", userID, query)
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

// cleanUpLoop runs in the background to remove expired entries.
func (c *AICache) cleanUpLoop() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.expiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// DefaultAICache is a global instance with a 10-minute TTL.
var DefaultAICache = NewAICache(10 * time.Minute)
