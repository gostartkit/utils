package utils

import (
	"container/list"
	"sync"
)

// CacheItem stores the key-value pair for the cache
type CacheItem struct {
	key   string
	value any
}

// LRUCache is a high-performance LRU cache
type LRUCache struct {
	items    map[string]*list.Element
	list     *list.List
	mu       sync.RWMutex
	capacity int
}

// NewLRUCache creates a new LRU cache with pre-allocated capacity
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 1 {
		capacity = 16
	}
	return &LRUCache{
		items:    make(map[string]*list.Element, capacity+1), // +1 to avoid immediate resize
		list:     list.New(),
		capacity: capacity,
	}
}

// Set adds or updates a key-value pair in the cache
func (c *LRUCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing key
	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*CacheItem).value = value
		return
	}

	// Evict LRU item if capacity is reached
	if c.list.Len() >= c.capacity {
		if lru := c.list.Back(); lru != nil {
			lruItem := lru.Value.(*CacheItem)
			delete(c.items, lruItem.key)
			c.list.Remove(lru)
		}
	}

	// Add new item
	item := &CacheItem{key: key, value: value}
	elem := c.list.PushFront(item)
	c.items[key] = elem
}

// Get retrieves a value from the cache
// Returns: value, exists
func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.RLock()

	if _, ok := c.items[key]; !ok {
		c.mu.RUnlock()
		return nil, false
	}

	// Move to front requires write lock
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[key]; ok { // Re-check in case of concurrent modification
		c.list.MoveToFront(elem)
		return elem.Value.(*CacheItem).value, true
	}
	return nil, false
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		delete(c.items, key)
		c.list.Remove(elem)
	}
}

// Len returns the current number of items in the cache
func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.list.Len()
}
