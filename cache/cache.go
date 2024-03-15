package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	capacity int
	data    map[string]*list.Element
	access *list.List
	mutex sync.Mutex
}

type CacheEntry struct {
	key string
	value interface{}
	timestamp time.Time
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		data: make(map[string]*list.Element),
		access: list.New(),
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if elem, exists := c.data[key]; exists {
		entry := elem.Value.(*CacheEntry)
		entry.value = value
		entry.timestamp = time.Now().Add(expiration)
		c.access.MoveToFront(elem)
		return
	}

	entry := &CacheEntry{
		key: key,
		value: value,
		timestamp: time.Now().Add(expiration),
	}

	elem := c.access.PushFront(entry)
	c.data[key] = elem

	if len(c.data) > c.capacity {
		c.evict()
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if elem, exists := c.data[key]; exists {
        // Move entry to front of access list (since it was accessed)
        c.access.MoveToFront(elem)
        // Return value associated with key
        return elem.Value.(*CacheEntry).value, true
    }

    // Key not found in cache
    return nil, false
}


func (c *Cache) evict() {
    if elem := c.access.Back(); elem != nil {
        // Remove least recently used entry from access list and data map
        entry := c.access.Remove(elem).(*CacheEntry)
        delete(c.data, entry.key)
    }
}