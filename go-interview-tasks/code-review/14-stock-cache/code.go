//go:build !solution

// Задача на code review (уровень: Lamoda).
// Кеш остатков SKU (Redis/in-memory) перед каталогом — типичный highload-паттерн Lamoda.
package stockcache

import "time"

type Loader func(sku string) (int, error)

type entry struct {
	qty       int
	expiresAt time.Time
}

type Cache struct {
	ttl  time.Duration
	data map[string]entry
}

func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, data: make(map[string]entry)}
}

func (c *Cache) Get(sku string, load Loader) (int, error) {
	if e, ok := c.data[sku]; ok && time.Now().Before(e.expiresAt) {
		return e.qty, nil
	}

	qty, err := load(sku)
	if err != nil {
		return 0, err
	}

	c.data[sku] = entry{qty: qty, expiresAt: time.Now().Add(c.ttl)}
	return qty, nil
}
