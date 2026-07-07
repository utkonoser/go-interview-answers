# Magnit Tech: atomic в Go

## M1. Как работают atomic операции под капотом в Go?

### Уровень CPU

Атомарность — **одна неделимая машинная инструкция**, а не «мьютекс спрятанный в int». На x86/arm64:

- `LOCK CMPXCHG` — compare-and-swap (CAS): «если значение всё ещё X — запиши Y»
- `XADD` — атомарное сложение
- `SWAP` — атомарная замена

`atomic.AddInt64(&v, 1)` компилируется в такую инструкцию. Другой поток не увидит промежуточное состояние «прочитал 5 → прибавил 1 → ещё не записал».

### Почему `count++` не атомарен

В Go `count++` — три шага: **load → add → store**. Две горутины могут обе прочитать `10`, обе записать `11` — один инкремент потерян. Это **data race** (`go test -race`).

### Память и happens-before

`sync/atomic` даёт гарантию видимости между горутинами на **одной** переменной: `Store` в одной горутине «виден» `Load` в другой. Обычные read/write без atomic компилятор и CPU могут переупорядочить — для флага «инициализация завершена» нужен `atomic` или mutex.

### API в Go

Низкоуровневые: `atomic.LoadInt64`, `StoreInt64`, `AddInt64`, `CompareAndSwapInt64`, `SwapInt64`, `uintptr`, `unsafe.Pointer`.

Обёртки (Go 1.19+): `atomic.Int64`, `atomic.Bool`, `atomic.Pointer[T]`, `atomic.Value`.

### Ограничения

Atomic защищает **одну ячейку**. Инвариант на двух полях (`reserved <= quantity`) одним atomic не обеспечить — нужен mutex или транзакция в БД.

Подходит: счётчики метрик, флаг shutdown, lock-free head списка, token bucket **в одном процессе**.

### Atomic vs Mutex

| | Atomic | Mutex |
|---|--------|-------|
| Скорость | Быстрее на простых ops | Дороже (contention, futex) |
| Сложность | Только простые операции | Любой инвариант под lock |
| Ошибки | Легко сломать логику CAS-loop | Проще reasoning |

На собесе Magnit часто спрашивают: «почему не `int` в счётчике RPS» — см. [06-rate-limiter.md](06-rate-limiter.md).
