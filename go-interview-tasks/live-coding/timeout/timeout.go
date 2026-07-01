// Пакет timeout — обёртка с дедлайном для долгой операции.
//
// Тяжёлую работу запускаем в горутине, ждём результат или отмену по context.
package timeout

import (
	"context"
	"time"
)

// Do выполняет fn с таймаутом d. При превышении — context.DeadlineExceeded.
func Do(fn func() error, d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel() // освобождаем таймер, даже если fn уже вернулась

	done := make(chan error, 1) // буфер 1 — горутина не зависнет после return
	go func() {
		done <- fn() // работа в отдельной горутине, чтобы select мог прервать ожидание
	}()

	select {
	case err := <-done:
		return err // успели до дедлайна
	case <-ctx.Done():
		return ctx.Err() // таймаут или отмена контекста
	}
}
