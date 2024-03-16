package main

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
