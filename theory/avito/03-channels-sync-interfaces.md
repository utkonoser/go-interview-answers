# Avito: каналы, sync, интерфейсы, память, тесты

## A3. Каналы

**Basic:** CSP-примитив для передачи значений и синхронизации; unbuffered = rendezvous; buffered = очередь фикс. размера; `select` мультиплексирует каналы.

**Advanced:**
- Запись в **closed** channel → panic; чтение — zero value + `ok=false`.
- Запись/чтение в **nil** channel блокируется навсегда (иногда намеренно в select).
- Каналы vs mutex: каналы — ownership transfer; mutex — shared state; «share memory by communicating».

---

## A3b. Примитивы синхронизации

| Примитив | Когда |
|----------|--------|
| `Mutex` / `RWMutex` | Инвариант на shared struct |
| `atomic` | Счётчики, flags, CAS на одной ячейке |
| `WaitGroup` | Дождаться N goroutines |
| `Once` | Однократная инициализация |
| `Cond` | Wait/signal на condition |
| channel | Pipeline, fan-in/out, cancellation signal |

Проблемы: **data race** (`go test -race`), **deadlock** (циклическое ожидание locks/channels).

---

## A3c. Интерфейсы

**Basic:** неявная реализация — тип с нужными методами; `interface{}` / `any` — любой тип.

**Advanced:** iface = `(type, data word)`; dynamic dispatch через itable.

**Expert — nil interface trap:**

```go
var p *Person = nil
var i interface{} = p
fmt.Println(i == nil) // false — type word != nil
```

---

## A3d. Память: slice и map

**len vs cap (slice):** len — видимые элементы; cap — ёмкость backing array от offset slice header.

**Рост slice:** при `append` beyond cap — новый array (~×2 до 1024, потом ~×1.25), copy, новый header.

**Map:** hash table; рост при load factor; iteration order random.

**GC (кратко):** mark reachable from roots (stacks, globals), sweep unmarked; STW только на короткие фазы; GOGC tuning.

---

## A3e. Тесты

- Виды: unit, integration, e2e, benchmark, fuzz.
- `go test`: `-run`, `-bench`, `-race`, `-cover`, `-count=1`, `-short`.
- Мocks: интерфейс + fake / testify/mock / gomock — подменить IO boundary.

---

## A3f. Профилирование

- **pprof:** CPU, heap, goroutine, mutex, block profiles.
- **trace:** timeline goroutines, GC, syscalls.
- CPU hot path → alloc reduction; mutex profile → contention; heap inuse → leak или oversized cache.

Включение pprof — sampling overhead, обычно приемлем в staging / short window in prod (с осторожностью).
