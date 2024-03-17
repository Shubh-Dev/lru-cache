package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	capacity int
	data     map[string]*list.Element
	access   *list.List
	mutex    sync.Mutex
}

type CacheEntry struct {
	key        string
	value      interface{}
	expiration time.Duration
	timestamp  time.Time
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		data:     make(map[string]*list.Element),
		access:   list.New(),
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.data[key]; exists {
		entry := elem.Value.(*CacheEntry)
		entry.value = value
		entry.expiration = expiration
		entry.timestamp = time.Now().Add(expiration)
		c.access.MoveToFront(elem)
		return
	}

	entry := &CacheEntry{
		key:        key,
		value:      value,
		expiration: expiration,
		timestamp:  time.Now(),
	}

	elem := c.access.PushFront(entry)
	c.data[key] = elem

	if len(c.data) > c.capacity {
		c.evict()
	}
}

func (c *Cache) Get(key string) (interface{}, time.Time, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.data[key]; exists {
		// Move entry to front of access list (since it was accessed)
		c.access.MoveToFront(elem)
		cacheEntry := elem.Value.(*CacheEntry)
		expirationTime := cacheEntry.timestamp.Add(cacheEntry.expiration)
		if time.Now().After(expirationTime) {
			c.access.Remove(elem)
			delete(c.data, key)
			return nil, expirationTime, false
		}

		return cacheEntry.value, expirationTime, true

	}
	return nil, time.Time{}, false
}

func (c *Cache) GetAllCache() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cacheContent := make(map[string]interface{})
	currentTime := time.Now()

	for key, elem := range c.data {
		cacheEntry := elem.Value.(*CacheEntry)
		expirationTime := cacheEntry.timestamp.Add(cacheEntry.expiration)

		if currentTime.After(expirationTime) {
			c.access.Remove(elem)
			delete(c.data, key)
		} else {
			cacheContent[key] = elem.Value.(*CacheEntry).value
		}
	}
	return cacheContent
}

func (c *Cache) evict() {
	if elem := c.access.Back(); elem != nil {
		entry := c.access.Remove(elem).(*CacheEntry)
		delete(c.data, entry.key)
	}
}
