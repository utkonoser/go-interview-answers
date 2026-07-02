// Пакет semaphore — ограничение числа одновременных операций.
//
// Отличается от worker pool: семафор лимитирует fan-out («не больше N HTTP-запросов сразу»),
// а не размер постоянной очереди воркеров.
package semaphore

import (
	"context"
	"sync"
	"time"
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

// FetchURLsParallel — пример: общий дедлайн + fan-out, но не больше maxConcurrent запросов сразу.
func FetchURLsParallel(parent context.Context, urls []string, maxConcurrent int) error {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	sem := New(maxConcurrent)
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			if err := sem.Acquire(ctx); err != nil {
				return // ctx отменён — слот не заняли, выходим
			}
			defer sem.Release()

			_ = fetchURL(ctx, url) // имитация HTTP к внешнему API
		}(url)
	}

	wg.Wait()
	return ctx.Err() // nil, если уложились в дедлайн
}

func fetchURL(ctx context.Context, url string) error {
	timer := time.NewTimer(50 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		_ = url
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
