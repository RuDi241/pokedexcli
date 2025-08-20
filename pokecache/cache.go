// Package pokecache: Provides a cache for pokedex app
package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mutex    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
		mutex:    sync.Mutex{},
	}
	go cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return []byte{}, ok
	}
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mutex.Lock()
		for k, v := range c.data {
			if time.Since(v.createdAt) > c.interval {
				delete(c.data, k)
			}
		}
		c.mutex.Unlock()
	}
}
