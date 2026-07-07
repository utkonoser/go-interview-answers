# Avito: горутины и runtime

## A2. Что такое горутины и как они соотносятся с потоками ОС?

**Basic (ожидание на собесе):**
- Легковесные «зелёные потоки», планируются **runtime Go**, не 1:1 с OS thread.
- Модель **M:N** — много goroutines на меньшее число OS threads (`GOMAXPROCS`).
- Переключение в user space дешевле, чем context switch потока (~KB стека goroutine vs ~MB на OS thread).

**Advanced:**
- Goroutine блокируется на syscall/network — runtime может отцепить её от P и занять M другой goroutine.
- Стек goroutine растёт динамически (от ~2 KB, copy-on-grow).
- Утечки goroutine: забыли выйти из select, нет consumer у channel, WaitGroup mismatch — искать через pprof goroutine + `-race`.

**Expert:**
- **G-M-P model:** G (goroutine), M (OS thread), P (logical processor с run queue).
- Локальная очередь P + global queue + work stealing.
- Sysmon будит goroutines после network poll / timer.

---

## A2b. Чем занимается runtime?

- Планирование goroutines (scheduler).
- Memory allocator + **GC** (tri-color concurrent mark-sweep).
- Stack management.
- **Netpoll** integration (epoll/kqueue/IOCP).
- `panic`/`recover`, defer, cgo bridge.

---

## A2c. Как Go держит много сетевых соединений?

**Basic:** **netpoller** — неблокирующие сокеты; goroutine «спит», не занимая OS thread, пока нет данных.

**Advanced:**
- Linux: **epoll** (edge/level triggered), BSD: kqueue, Windows: IOCP.
- Преимущество epoll над select/poll: O(1) по числу fd, без линейного scan.
- Регистрируем интересующие events: read/write/error.

**Expert:** edge-triggered — событие один раз до следующего `read`; level-triggered — пока буфер не пуст. Go runtime абстрагирует это в `internal/poll`.
