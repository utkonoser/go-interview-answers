//go:build solution

package stockcache

import (
	"sync"
	"time"
)

type Loader func(sku string) (int, error)

type entry struct {
	qty       int
	expiresAt time.Time
}

type Cache struct {
	ttl      time.Duration
	mu       sync.Mutex
	data     map[string]entry
	inflight map[string]chan struct{}
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		ttl:      ttl,
		data:     make(map[string]entry),
		inflight: make(map[string]chan struct{}),
	}
}

// fix: data race на map; fix: thundering herd — coalescing через inflight channel.
func (c *Cache) Get(sku string, load Loader) (int, error) {
	c.mu.Lock()
	if e, ok := c.data[sku]; ok && time.Now().Before(e.expiresAt) {
		c.mu.Unlock()
		return e.qty, nil
	}
	if ch, ok := c.inflight[sku]; ok {
		c.mu.Unlock()
		<-ch
		c.mu.Lock()
		e := c.data[sku]
		c.mu.Unlock()
		return e.qty, nil
	}

	done := make(chan struct{})
	c.inflight[sku] = done
	c.mu.Unlock()

	qty, err := load(sku)

	c.mu.Lock()
	if err == nil {
		c.data[sku] = entry{qty: qty, expiresAt: time.Now().Add(c.ttl)}
	}
	delete(c.inflight, sku)
	close(done)
	c.mu.Unlock()

	return qty, err
}
