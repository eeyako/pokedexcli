package internal

import (
	"sync"
	"time"
)

type Cache struct {
	CacheEntry map[string]CacheEntry
	sync.Mutex
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       *[]byte
}

func NewCache(cache *Cache, key string, val *[]byte, interval time.Duration) {
	cache.ReapLoop(interval)
	cache.Add(key, val)
}

func (c *Cache) Add(key string, val *[]byte) {
	for k := range c.CacheEntry {
		if key == k {
			return
		}
	}
	newEntry := CacheEntry{CreatedAt: time.Now(), Val: val}
	if c.CacheEntry == nil {
		c.CacheEntry = make(map[string]CacheEntry)
	}
	c.CacheEntry[key] = newEntry
}

func (c *Cache) Get(key string) (*[]byte, bool) {
	for k, cacheEntry := range c.CacheEntry {
		if key == k {
			return cacheEntry.Val, true
		}
	}

	return &[]byte{}, false
}

func (c *Cache) ReapLoop(interval time.Duration) {
	currentTime := time.Now()
	cutoff := currentTime.Add(-interval)
	for key, cacheEntry := range c.CacheEntry {
		if cacheEntry.CreatedAt.Before(cutoff) {
			delete(c.CacheEntry, key)
		}
	}
}
