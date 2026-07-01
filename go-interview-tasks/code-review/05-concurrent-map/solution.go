//go:build solution

package concurrentmap

import "sync"

// fix: map не потокобезопасна — concurrent read/write паникует или даёт data race.
type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}
