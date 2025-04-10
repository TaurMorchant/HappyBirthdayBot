package cache

import (
	"fmt"
	gocache "github.com/patrickmn/go-cache"
	"log"
	"time"
)

type Cache[K any, V interface{}] struct {
	gocache.Cache
}

func New[K any, V interface{}](defaultExpiration, cleanupInterval time.Duration) *Cache[K, V] {
	result := Cache[K, V]{Cache: *gocache.New(defaultExpiration, cleanupInterval)}
	return &result
}

func (c *Cache[K, V]) Add(key K, value V) {
	c.Cache.Set(fmt.Sprintf("%v", key), value, gocache.DefaultExpiration)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	log.Println("[TRACE] All elements in cache: ", c.Items())
	result, ok := c.Cache.Get(fmt.Sprintf("%v", key))
	if ok {
		return result.(V), ok
	} else {
		var zero V
		return zero, ok
	}
}

func (c *Cache[K, V]) Delete(key K) {
	c.Cache.Delete(fmt.Sprintf("%v", key))
}
