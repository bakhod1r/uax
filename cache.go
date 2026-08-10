package uax

import (
	"container/list"
	"sync"
)

// cacheItem stores the value and the pointer to the list element.
type cacheItem struct {
	key   string
	value *Client
}

// Cache is a simple, thread-safe LRU cache for parsed Clients.
type Cache struct {
	maxEntries int
	ll         *list.List
	cache      map[string]*list.Element
	mu         sync.RWMutex
}

// NewCache creates a new LRU cache.
func NewCache(maxEntries int) *Cache {
	if maxEntries <= 0 {
		maxEntries = 1000 // default
	}
	return &Cache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (*Client, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*cacheItem).value, true
	}
	return nil, false
}

// Set adds a value to the cache.
func (c *Cache) Set(key string, value *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		ele.Value.(*cacheItem).value = value
		return
	}

	ele := c.ll.PushFront(&cacheItem{key, value})
	c.cache[key] = ele
	if c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*cacheItem)
		delete(c.cache, kv.key)
	}
}
