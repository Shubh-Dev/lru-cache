package cache_test

import (
	"testing"
	"time"

	"github.com/Shubh-Dev/lru-cache/cache"
)

func TestSetAndGet(t *testing.T) {
	expected := "value"
	cacheInstance := cache.NewCache(10)
	cacheInstance.Set("key", expected, 1*time.Second)
	actual, _, found := cacheInstance.Get("key")
	if !found || actual != expected {
		t.Errorf("Set and Get failed: expected %s, got %s", expected, actual)
	}
}

func TestExpiry(t *testing.T) {
	cacheInstance := cache.NewCache(10)

	// Set an entry with a short expiration time
	cacheInstance.Set("key1", "value1", 1*time.Second)

	// Wait for the expiration time to elapse
	time.Sleep(2 * time.Second)

	// Try to retrieve the expired entry
	_, _, found := cacheInstance.Get("key1")
	if found {
		t.Error("Expected entry to be expired, but it was found in the cache")
	}

	// Set a new entry with a longer expiration time
	cacheInstance.Set("key2", "value2", 5*time.Second)

	// Retrieve the new entry before it expires
	value, _, found := cacheInstance.Get("key2")
	if !found {
		t.Error("Expected entry to be found in the cache, but it was not found")
	}
	if value != "value2" {
		t.Errorf("Expected value 'value2', got '%s'", value)
	}

	// Wait for the new entry to expire
	time.Sleep(6 * time.Second)

	// Try to retrieve the expired entry
	_, _, found = cacheInstance.Get("key2")
	if found {
		t.Error("Expected entry to be expired, but it was found in the cache")
	}
}
