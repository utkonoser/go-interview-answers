// Пакет errgroup — параллельный запуск задач с отменой при первой ошибке.
//
// Обёртка над golang.org/x/sync/errgroup: типичный паттерн на собесах Avito/Ozon/Яндекс
// («сходи в 3 сервиса параллельно, при ошибке одного — отмени остальные»).
package errgroup

import (
	"context"

	"golang.org/x/sync/errgroup"
)

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
