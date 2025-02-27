package pokecache

import (
	"sync"
	"time"
)

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheEntries: make(map[string]cacheEntry),
	}
	cache.reapLoop(interval)
	return cache
}

type Cache struct {
	mu           sync.Mutex
	cacheEntries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte // raw data we're caching
}

func (c *Cache) Add(key string, value []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	// Get entry from the cache
	// If the entry is found, return true
	// Otherwise return false
	result, ok := c.cacheEntries[key]
	if ok {
		return result.val, ok
	}
	return nil, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	// Each time an interval (the time.Duration passed to NewCache) passes it should remove any entries that are older than the interval
	// Example: If the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed
	// Create a new ticker that sends a value on its channel (ticker.C) at regular intervals (which we specify when we create a new ticker)
	ticker := time.NewTicker(interval)
	// Anonymous Go routine that will run in the background
	go func() {
		// Runs indefinitely, waiting for each tick
		for {
			// Wait until next tick
			<-ticker.C
			// Lock the mutext before accessing the map - ensures thread safety
			c.mu.Lock()
			// Get current time
			now := time.Now()
			// Check entry in the map
			for key, entry := range c.cacheEntries {
				// If the entry is older than the interval, delete it
				if now.Sub(entry.createdAt) > interval {
					delete(c.cacheEntries, key)
				}
			}
			// Unlock mutex
			c.mu.Unlock()
		}
	}()
}
