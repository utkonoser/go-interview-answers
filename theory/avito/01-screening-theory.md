# Avito: скрининг — теория

## A1. Go: слайс, map, set, каналы, context, panic

**Функция как параметр?** Да, функции — first-class values (функции высшего порядка).

**Слайс vs массив:** массив — фиксированный размер, value type; слайс — descriptor `(ptr, len, cap)` поверх массива, reference semantics.

**Map lookup:** в среднем **O(1)**, worst case O(n) при коллизиях/ребалансе.

**Set в Go:** нет в stdlib; `map[T]struct{}` — нулевой value size для `struct{}`.

**Каналы:** буферизованные (`make(chan T, n)`) и небуферизованные (синхронная handoff).

**context.Context:** отмена, дедлайны, request-scoped values; прокидывать по границам API.

**Exceptions:** нет; для программных ошибок — `error`, для невосстановимого — `panic` + `recover` на границе (HTTP middleware, worker top-level).

**Mutex:** mutual exclusion для shared memory; `RWMutex` — много readers или один writer.

---

## A2. БД: индексы, B-Tree, транзакции, WHERE/HAVING, триггеры

**Зачем индексы:** ускорить поиск/сортировку/join; примеры — B-Tree (PostgreSQL default), Hash, GiST/GIN, R-Tree.

**B-Tree быстрее full scan:** сбалансированное дерево, высота O(log n), страницы на диске — меньше random I/O.

**Транзакция:** несколько операций атомарно; **ACID** — Atomicity, Consistency, Isolation, Durability.

**WHERE vs HAVING:** `WHERE` фильтрует строки **до** `GROUP BY`; `HAVING` — группы **после** агрегации.

**Агрегаты:** `COUNT`, `SUM`, `AVG`, `MIN`, `MAX` — фильтр по результату агрегации в `HAVING`.

**Триггер:** процедура на `INSERT/UPDATE/DELETE`; audit, denormalization, constraints — осторожно с latency и отладкой.

---

## A3. Сеть, HTTP, OS, контейнеры

**TCP vs UDP:** TCP — connection, порядок, retransmit; UDP — без гарантий, меньше overhead.

**HTTP request:** Request line (method, path, version) → Headers → blank line → Body.

**HTTPS vs HTTP:** TLS поверх TCP — шифрование, аутентификация сервера (сертификат).

**HTTP timeouts:** client timeout / context deadline; подбирать по p99 latency + запас, не «навсегда».

**Поток vs процесс:** процесс — отдельное address space; потоки одного процесса — общая память.

**kill:** `kill PID`, `killall name`, `pkill pattern`.

**Выделение 1 KB:** allocator runtime → возможно новый span от OS; при нехватке — GC, затем OOM kill.

**Pod:** группа контейнеров, shared network namespace и volumes в k8s.

**Docker vs VM:** контейнеры делят kernel хоста — легче, быстрее старт; VM — полная guest OS.
