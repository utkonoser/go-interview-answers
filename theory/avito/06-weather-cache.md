# Avito: прогноз погоды 10k RPS

## A6. Прогноз погоды при 10k RPS

- `aiWeatherForecast()` — ~1 сек CPU/IO (NN).
- HTTP `/weather` — **10k RPS**.
- Нельзя вызывать NN на каждый request.

### Антипаттерн

```go
http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, `{"temperature":%d}`, aiWeatherForecast()) // 1 RPS effective
})
```

---

## Решение: кеш + фоновое обновление

```go
type cache struct {
    mu          sync.RWMutex
    temperature int
}

var c cache

func refresh(ctx context.Context) {
    for {
        v := aiWeatherForecast()
        c.mu.Lock()
        c.temperature = v
        c.mu.Unlock()
        select {
        case <-ctx.Done():
            return
        case <-time.After(time.Second):
        }
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    c.mu.RLock()
    t := c.temperature
    c.mu.RUnlock()
    fmt.Fprintf(w, `{"temperature":%d}`, t)
}
```

Handler — **O(1)**, read lock; refresh — один goroutine или worker pool.

---

## С `cityId`

- Ключ кеша: `cityId` → `map[int]entry{temperature, updatedAt}`.
- Структура: `map` + `RWMutex` или `sync.Map` для read-heavy; **не vector** — cityId sparse.
- Per-key **singleflight** или `Semaphore(1)` — thundering herd при miss не запускает 1000 NN calls.

```go
import "golang.org/x/sync/singleflight"

var g singleflight.Group

func getWeather(cityID int) (int, error) {
    v, err, _ := g.Do(strconv.Itoa(cityID), func() (any, error) {
        return aiWeatherForecast(cityID)
    })
    return v.(int), err
}
```

---

## Нештатные ситуации кеша

| Проблема | Решение |
|----------|---------|
| Cold start | warmup top-N cities при старте |
| TTL expiry thundering herd | jitter TTL, singleflight, stale-while-revalidate |
| Stale data | TTL + background refresh до expiry |
| Race на map | RWMutex / sync.Map / sharded locks |
| High write contention | отдельный refresh goroutine, readers только читают snapshot |

**Таймаут на NN:** context 1–2s; при timeout — last known value или 503.

**HTTP:** 200 OK, 400 bad cityId, 500 upstream fail — не маскировать 500 как 200.

**Expert follow-ups:** double buffering (swap `storage`/`oldStorage`); periodic top-N refresh; `ConcurrentDictionary` аналог в других языках.

См. `go-interview-tasks/live-coding/singleflight/`, `ratelimit/`.
