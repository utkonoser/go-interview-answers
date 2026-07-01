//go:build !solution

// Задача на code review: in-memory кеш для HTTP handler'ов.
package concurrentmap

// Cache хранит строковые значения по ключу.
type Cache struct {
	data map[string]string
}

func New() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Set(key, value string) {
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	v, ok := c.data[key]
	return v, ok
}
