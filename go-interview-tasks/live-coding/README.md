# Live-coding: типовые паттерны на Go

Эталонные реализации для live-coding (Avito, Ozon, WB, Яндекс): LRU, worker pool, rate limiter, graceful shutdown.

| # | Пакет | Тема |
|---|--------|------|
| 1 | `lru/` | LRU cache (map + двусвязный список) |
| 2 | `counter/` | Thread-safe счётчик (`atomic` / `mutex`) |
| 3 | `workerpool/` | Worker pool |
| 4 | `ratelimit/` | Token bucket rate limiter |
| 5 | `producerconsumer/` | Producer-consumer |
| 6 | `merge/` | Merge каналов |
| 7 | `timeout/` | Timeout через `context` |
| 8 | `shutdown/` | Graceful HTTP shutdown |
| 9 | `errgroup/` | Параллельные задачи с отменой при ошибке |
| 10 | `singleflight/` | Один запрос на ключ (thundering herd) |
| 11 | `semaphore/` | Лимит одновременных операций |
| 12 | `retry/` | Exponential backoff + jitter |

```bash
cd go-interview-tasks
go test ./live-coding/tests/...
```
