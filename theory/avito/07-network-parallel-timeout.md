# Avito: сеть, параллельные запросы, timeout

## A7. IO-bound: 10 000 последовательных «сетевых» вызовов

Последовательно: ~10 s (1 ms × 10k). Параллельно:

```go
var wg sync.WaitGroup
var count atomic.Int64
for i := 0; i < numRequests; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        networkRequest()
        count.Add(1)
    }()
}
wg.Wait()
```

**Ожидание:** ~1–2 ms wall time (если IO parallel). **Но:** 10k goroutines на **реальный** HTTP — risk: FD limit, upstream rate limit, memory.

**Ограничение:** semaphore / worker pool (например 100 concurrent).

---

## A7b. Параллельные HTTP status codes

```go
func fetchStatus(ctx context.Context, url string) error {
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    io.Copy(io.Discard, resp.Body)
    fmt.Println(url, resp.Status)
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    g, ctx := errgroup.WithContext(ctx)
    for _, u := range urls {
        u := u
        g.Go(func() error { return fetchStatus(ctx, u) })
    }
    _ = g.Wait()
}
```

**Ошибки кандидатов:** `wg.Add` внутри goroutine; sequential fetch; нет drain body (reuse connection).

**Follow-ups:**
- 500 URLs — тот же паттерн, pool ограничивает fan-out.
- Max 2s — `context.WithTimeout`.
- Graceful shutdown — `signal.Notify`, `server.Shutdown(ctx)`.
- HTTP/1.1 keep-alive — `Transport.MaxIdleConns`, `IdleConnTimeout`.
- Debug connections — `ss`, `lsof -p PID`.
- Circuit breaker при flaky upstream.

---

## A7c. predictableFunc — timeout обёртка

```go
func predictableFunc(ctx context.Context) (int64, error) {
    ctx, cancel := context.WithTimeout(ctx, time.Second)
    defer cancel()

    start := time.Now()
    defer func() { log.Println("elapsed:", time.Since(start)) }()

    ch := make(chan int64, 1) // buffered — иначе leak goroutine при timeout
    go func() {
        ch <- unpredictableFunc()
    }()

    select {
    case v := <-ch:
        return v, nil
    case <-ctx.Done():
        return 0, ctx.Err()
    }
}
```

**Критично:**
- **Buffered channel** — иначе goroutine зависнет после timeout.
- **`defer cancel()`** — освободить timer.
- `defer` с `time.Since(start)` — захват `start` в closure, не `time.Since(time.Now())` в defer args.

**select с несколькими ready channels** — choice **pseudo-random**, не полагаться на порядок.

**Лучше (если можно):** `unpredictableFunc(ctx context.Context)` — native cancellation.

См. `go-interview-tasks/live-coding/timeout/`, `shutdown/`.
