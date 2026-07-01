// Пакет semaphore — ограничение числа одновременных операций.
//
// Отличается от worker pool: семафор лимитирует fan-out («не больше N HTTP-запросов сразу»),
// а не размер постоянной очереди воркеров.
package semaphore

import (
	"context"
)

// Semaphore — счётчик слотов через буферизованный канал.
type Semaphore struct {
	slots chan struct{}
}

// New создаёт семафор на n одновременных слотов.
func New(n int) *Semaphore {
	if n < 1 {
		n = 1
	}
	return &Semaphore{slots: make(chan struct{}, n)}
}

// Acquire занимает слот. Блокируется, пока слот не освободится или ctx не отменён.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release освобождает слот. Вызывай в defer после успешного Acquire.
func (s *Semaphore) Release() {
	<-s.slots
}
