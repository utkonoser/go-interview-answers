// Пакет ratelimit — ограничение частоты запросов (token bucket).
//
// «Жетоны» лежат в буферизованном канале; фоновая горутина периодически их пополняет.
package ratelimit

import "time"

// Limiter — token bucket: Allow забирает жетон, если есть.
type Limiter struct {
	tokens chan struct{} // буфер = максимум жетонов в bucket
	ticker *time.Ticker  // равномерно пополняет bucket
}

// New создаёт лимитер: rate запросов за интервал per.
func New(rate int, per time.Duration) *Limiter {
	if rate < 1 {
		rate = 1
	}
	rl := &Limiter{
		tokens: make(chan struct{}, rate),
		ticker: time.NewTicker(per / time.Duration(rate)),
	}
	// стартуем с полным bucket — иначе первые запросы всегда отклонятся
	for range rate {
		rl.tokens <- struct{}{}
	}
	// фоновое пополнение: один жетон каждые per/rate
	go func() {
		for range rl.ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default: // bucket полон — лишний жетон выбрасываем
			}
		}
	}()
	return rl
}

// Allow возвращает true, если запрос разрешён (жетон был в bucket).
func (rl *Limiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default: // жетонов нет — лимит исчерпан
		return false
	}
}

// Stop останавливает тикер пополнения (вызывай при завершении, иначе утечка горутины).
func (rl *Limiter) Stop() {
	rl.ticker.Stop()
}
