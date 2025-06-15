// https://www.boot.dev/lessons/f2ba2b87-38fe-467c-abd7-14716f955169

// Mutex:
// https://blog.boot.dev/golang/golang-mutex/
// https://gobyexample.com/mutexes

package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

// You'll probably want to expose a NewCache() function that creates a new cache with a configurable interval (time.Duration).
func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		entries:  map[string]cacheEntry{},
		interval: interval,
	}

	go newCache.reapLoop()

	return &newCache
}

// Create a cache.Add() method that adds a new entry to the cache.
// It should take a key (a string) and a val (a []byte).
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()

	newEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = newEntry

	c.mu.Unlock()
}

// Create a cache.Get() method that gets an entry from the cache.
// It should take a key (a string) and return a []byte and a bool.
// The bool should be true if the entry was found and false if it wasn't.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()

	entry, ok := c.entries[key]

	c.mu.Unlock()

	return entry.val, ok
}

// Create a cache.reapLoop() method that is called when the cache is created (by the NewCache function).
// Each time an interval (the time.Duration passed to NewCache) passes
// it should remove any entries that are older than the interval.
// This makes sure that the cache doesn't grow too large over time.
// For example, if the interval is 5 seconds, and an entry was added 7 seconds ago, that entry should be removed.
func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)

	for {
		// wait for a tick
		<-ticker.C

		c.mu.Lock()

		for key := range c.entries {
			entry := c.entries[key]
			if entry.createdAt.Add(c.interval).Before(time.Now()) {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()
	}
}

// I used a time.Ticker to make this happen.
// https://pkg.go.dev/time#Ticker
