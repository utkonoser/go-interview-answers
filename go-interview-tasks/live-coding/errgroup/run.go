// Пакет errgroup — параллельный запуск задач с отменой при первой ошибке.
//
// Обёртка над golang.org/x/sync/errgroup: типичный паттерн на собесах Avito/Ozon/Яндекс
// («сходи в 3 сервиса параллельно, при ошибке одного — отмени остальные»).
package errgroup

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

// ParallelFetch — пример вызова: общий дедлайн + три «сервиса» параллельно.
// При ошибке или таймауте одного ctx отменяется — остальные должны выйти по ctx.Done().
func ParallelFetch(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, 3*time.Second)
	defer cancel()

	return RunParallel(ctx,
		func(ctx context.Context) error { return fetchProfile(ctx) },
		func(ctx context.Context) error { return fetchOrders(ctx) },
		func(ctx context.Context) error { return fetchBalance(ctx) },
	)
}

// fetchProfile имитирует HTTP/БД-вызов: ждём или отменяемся по ctx.
func fetchProfile(ctx context.Context) error {
	return waitOrCancel(ctx, 50*time.Millisecond)
}

func fetchOrders(ctx context.Context) error {
	return waitOrCancel(ctx, 100*time.Millisecond)
}

func fetchBalance(ctx context.Context) error {
	return waitOrCancel(ctx, 80*time.Millisecond)
}

func waitOrCancel(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err() // errgroup отменил ctx — выходим без лишней работы
	}
}

// RunParallel запускает fns параллельно в одной errgroup.
// При первой ошибке контекст отменяется — остальные задачи должны слушать ctx.Done().
func RunParallel(ctx context.Context, fns ...func(context.Context) error) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, fn := range fns {
		g.Go(func() error {
			return fn(ctx)
		})
	}

	return g.Wait()
}
