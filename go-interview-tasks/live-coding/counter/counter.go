// Пакет counter — потокобезопасные счётчики для конкурентного инкремента.
//
// Atomic — для простого +1 (быстрее, без блокировок).
// Mutex — когда нужна сложнее логика под одной блокировкой.
package counter

import (
	"sync"
	"sync/atomic"
)

// Atomic — счётчик на atomic-операциях CPU.
type Atomic struct {
	value int64
}

// Increment атомарно увеличивает значение на 1.
func (c *Atomic) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Value атомарно читает текущее значение.
func (c *Atomic) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Mutex — счётчик под sync.Mutex.
type Mutex struct {
	mu    sync.Mutex
	value int64
}

// Increment увеличивает значение под мьютексом.
func (c *Mutex) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // разблокируем при выходе из функции
	c.value++
}

// Value возвращает значение под мьютексом.
func (c *Mutex) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
