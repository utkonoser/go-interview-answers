// Пакет retry — повтор вызова с exponential backoff и jitter.
//
// Для нестабильных HTTP/БД: transient errors, 503, timeout. Пара к code-review/04-http-client.
package retry

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Config задаёт политику повторов.
type Config struct {
	MaxAttempts int           // всего попыток, включая первую
	BaseDelay   time.Duration // начальная пауза между попытками
	MaxDelay    time.Duration // потолок паузы
}

// DefaultConfig — разумные значения для собеса/прода.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    2 * time.Second,
	}
}

// Do вызывает fn до успеха или исчерпания попыток. Между попытками — backoff + jitter.
// Уважает отмену ctx во время ожидания.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 1
	}

	var err error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts {
			break
		}

		delay := backoff(cfg, attempt)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.Join(err, ctx.Err())
		case <-timer.C:
		}
	}
	return err
}

// backoff — экспоненциальная задержка с jitter, чтобы не бить в сервис синхронно.
func backoff(cfg Config, attempt int) time.Duration {
	delay := cfg.BaseDelay * time.Duration(1<<(attempt-1))
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	// jitter: от 50% до 150% от delay
	jitter := time.Duration(rand.Int63n(int64(delay))) - delay/2
	return delay + jitter
}
