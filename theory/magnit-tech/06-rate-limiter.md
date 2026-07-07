# Magnit Tech: rate limiter

## M6. Rate limiter на PostgreSQL и в Go на atomic (почему не `int`)

### Задача

Лимит **N запросов в минуту** на `api_key` при **нескольких репликах** API — счётчик в памяти одного pod не подходит.

---

### Вариант 1: PostgreSQL (fixed window)

```sql
CREATE TABLE rate_limits (
    api_key    TEXT NOT NULL,
    window_start TIMESTAMPTZ NOT NULL,
    count      INT NOT NULL DEFAULT 0,
    PRIMARY KEY (api_key, window_start)
);

-- В транзакции, window_start = date_trunc('minute', now())
INSERT INTO rate_limits (api_key, window_start, count)
VALUES ($1, $2, 1)
ON CONFLICT (api_key, window_start)
DO UPDATE SET count = rate_limits.count + 1
RETURNING count;
```

Если `count > N` → 429 Too Many Requests.

**Плюсы:** один источник правды для всех реплик.  
**Минусы:** нагрузка на БД на каждый запрос; для high RPS лучше Redis + TTL или sliding window в памяти gateway.

**Атомарность:** `INSERT ... ON CONFLICT DO UPDATE` — одна строка, row-level lock, гонки нет.

**С advisory lock** (альтернатива):

```sql
SELECT pg_advisory_xact_lock(hashtext($api_key));
-- read count, increment, check limit
```

Держит lock до конца транзакции.

---

### Вариант 2: Go in-process на `atomic` (одна реплика)

```go
type Limiter struct {
    tokens atomic.Int64
}

func (l *Limiter) Allow() bool {
    for {
        cur := l.tokens.Load()
        if cur <= 0 {
            return false
        }
        if l.tokens.CompareAndSwap(cur, cur-1) {
            return true
        }
    }
}
```

Пополнение: горутина с `time.Ticker` → `atomic.Store` / `Add`.

**Почему CAS-loop:** два `Load` + `Store` без CAS — lost update между горутинами.

---

### Почему обычный `int` не работает

```go
var tokens int

func allow() bool {
    if tokens <= 0 {
        return false
    }
    tokens--  // DATA RACE: не атомарно
    return true
}
```

Две горутины обе читают `tokens == 1`, обе проходят проверку, обе декрементят — лимит нарушен. `go test -race` покажет гонку.

Даже с `sync.Mutex` вокруг `int` — ок для одного процесса, но **не для кластера**.

---

### Сравнение

| Решение | Несколько реплик | Latency | Сложность |
|---------|------------------|---------|-----------|
| `atomic` in-memory | ✗ | Низкая | Низкая |
| `int` без sync | ✗ (и race) | — | Баг |
| PostgreSQL upsert | ✓ | Выше | Средняя |
| Redis INCR+EXPIRE | ✓ | Низкая | Средняя |
| API Gateway (nginx) | ✓ | Низкая | Инфра |

### На собесе Magnit

Ожидают: понимание **где** state (память vs БД vs Redis), **почему** `int++` race, **как** PG обеспечивает атомарный инкремент строки, когда БД — overkill (миллионы RPS → Redis/gateway).

См. также live-coding: `go-interview-tasks/live-coding/ratelimit/`, `atomic` в `counter/`.
