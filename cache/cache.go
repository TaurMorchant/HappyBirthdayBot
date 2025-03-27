package cache

import (
	"fmt"
	gocache "github.com/patrickmn/go-cache"
	"log"
	"time"
)

type Cache[K, V any] struct {
	gocache.Cache
}

func New[K, V any](defaultExpiration, cleanupInterval time.Duration) *Cache[K, V] {
	result := Cache[K, V]{Cache: *gocache.New(defaultExpiration, cleanupInterval)}
	return &result
}

func (c *Cache[K, V]) Add(key K, value V) {
	err := c.Cache.Add(fmt.Sprintf("%d", key), value, gocache.DefaultExpiration)
	if err != nil {
		log.Panicln("Cannot add element to cache: ", err)
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	log.Println("[TRACE] All elements in cache: ", c.Items())
	result, ok := c.Cache.Get(fmt.Sprintf("%d", key))
	return result.(V), ok
}

func (c *Cache[K, V]) Delete(key K) {
	c.Cache.Delete(fmt.Sprintf("%d", key))
}
