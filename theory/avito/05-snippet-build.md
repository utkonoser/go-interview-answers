# Avito: сборка сниппета

## A5. Сборка сниппета товара

Собрать сниппет товара: **описание** (сервис → `prettify`) + **цена в рублях** (сервис в USD → `priceToRub`). Два независимых IO-вызова.

### Базовое решение — параллельно через WaitGroup

```go
func BuildSnippet(itemID int) (Snippet, error) {
    var wg sync.WaitGroup
    wg.Add(2)

    var desc string
    var priceRub float64
    var errDesc, errPrice error

    go func() {
        defer wg.Done()
        raw, err := itemDescription(itemID)
        if err != nil {
            errDesc = err
            return
        }
        desc = prettify(raw)
    }()

    go func() {
        defer wg.Done()
        usd, err := itemPrice(itemID)
        if err != nil {
            errPrice = err
            return
        }
        priceRub = priceToRub(usd)
    }()

    wg.Wait()
    if errDesc != nil {
        return Snippet{}, errDesc
    }
    if errPrice != nil {
        return Snippet{}, errPrice
    }
    return Snippet{Description: desc, Price: priceRub}, nil
}
```

### Идиоматичнее — errgroup + context

```go
func BuildSnippet(ctx context.Context, itemID int) (Snippet, error) {
    g, ctx := errgroup.WithContext(ctx)
    var desc string
    var priceRub float64

    g.Go(func() error {
        raw, err := itemDescription(ctx, itemID)
        if err != nil {
            return err
        }
        desc = prettify(raw)
        return nil
    })
    g.Go(func() error {
        usd, err := itemPrice(ctx, itemID)
        if err != nil {
            return err
        }
        priceRub = priceToRub(usd)
        return nil
    })

    if err := g.Wait(); err != nil {
        return Snippet{}, err
    }
    return Snippet{Description: desc, Price: priceRub}, nil
}
```

Первая ошибка отменяет `ctx` — второй вызов может прерваться.

---

## Дополнительные вопросы (как на собесе)

**Пачка задач:** `errgroup` или worker pool с semaphore; ограничить parallelism.

**Ошибка одного сервиса:**
- fail whole snippet (500) — строгий контракт;
- partial response (цена без описания) — явно в API contract;
- fallback (USD вместо RUB) — только если продукт согласен;
- retry с backoff — для transient errors, idempotent GET.

**Таймауты:** `context.WithTimeout` на handler; per-upstream timeout ≤ общего; учитывать p99 dependency + margin.

**Threads vs goroutines:** OS threads тяжёлые; goroutines multiplexed на GOMAXPROCS; блокирующий IO не блокирует весь процесс.

**Scheduler:** goroutine park на channel/network/syscall; wake через netpoll/timer.

**Mutex:** always `defer Unlock`; RWMutex если read-heavy; panic в critical section — `defer` всё равно отпустит lock (если не recursive misuse).

**Atomic/CAS:** счётчики OK; составной invariant — mutex.

**Sync между машинами:** не shared memory — БД, Redis lock, message queue, idempotency keys.

См. также `go-interview-tasks/live-coding/errgroup/`.
