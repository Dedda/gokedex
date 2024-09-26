package gokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex   sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		mutex:   sync.RWMutex{},
	}
	ticker := time.NewTicker(interval)
	go c.reapLoop(interval, ticker.C)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.entries[key]; !ok {
		c.entries[key] = cacheEntry{
			createdAt: time.Now(),
			val:       val,
		}
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if e, ok := c.entries[key]; ok {
		return e.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration, channel <-chan time.Time) {
	for t := range channel {
		c.mutex.Lock()
		for k, v := range c.entries {
			if v.createdAt.Add(interval).Before(t) {
				delete(c.entries, k)
			}
		}
		c.mutex.Unlock()
	}
}
