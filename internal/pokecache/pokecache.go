package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu        sync.RWMutex
	entries   map[string]cacheEntry
	timeLimit time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(duration time.Duration) *Cache {
	cache := Cache{
		mu:        sync.RWMutex{},
		entries:   map[string]cacheEntry{},
		timeLimit: duration,
	}

	// By initiating a go routine inside the constructor,
	// I can keep track of the reapLoop timer
	go cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists := c.entries[key]
	if exists {
		return item.val, true
	}

	return nil, false
}

// Will continuously check for a tick
func (c *Cache) reapLoop() {
	// Ticker is initialized with a channel
	ticker := time.NewTicker(c.timeLimit)
	for {
		select {
		// When the duration is up, channel will receive a value
		case <-ticker.C:
			// Calls reap when ticker is up
			c.reap()
		}
	}
}

func (c *Cache) reap() {
	currentTime := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, cacheEntry := range c.entries {
		// Calculate the time allowed for each entry based on the time limit given
		expirationTime := cacheEntry.createdAt.Add(c.timeLimit)
		if currentTime.After(expirationTime) {
			delete(c.entries, key)
		}
	}
}
