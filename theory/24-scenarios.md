# Сценарии и troubleshooting

Вопросы формата «есть проблема — расскажи, как решать». На собесе ждут пошаговый план, а не угадывание одной причины.

## 200. Latency выросла после деплоя — с чего начать?

1. **Сопоставить время** — графики RPS, p50/p95/p99, error rate в Grafana; совпадает ли с деплоем/изменением конфига.
2. **Откатить или канареечный pod** — если явная корреляция, быстро проверить гипотезу rollback'ом.
3. **Трейсы** — один медленный запрос: какой span съел время (БД, внешний API, сериализация).
4. **pprof** — CPU profile, сравнить с прошлой версией.
5. **БД** — slow query log, `pg_stat_activity`, рост locks.
6. **Diff деплоя** — новая зависимость, убранный индекс, изменённый timeout, feature flag.
7. **Нагрузка** — может latency выросла из-за роста RPS, а не из-за кода.

## 201. Растёт число горутин (goroutine leak) — как найти?

1. Метрика `go_goroutines` в Prometheus — тренд вверх без спада после пика нагрузки.
2. **pprof goroutine**: `go tool pprof http://localhost:6060/debug/pprof/goroutine` — stack trace всех горутин, ищем одинаковые стеки (застряли на `chan receive`, `sync.Mutex.Lock`, `time.Sleep`).
3. Типичные причины: горутина ждёт закрытия канала, который никто не закроет; забытый `ctx` без cancel; worker без exit condition.
4. Локально: `goleak` в тестах.
5. Фикс: `context` с cancel, закрытие каналов отправителем, таймауты на блокирующие операции.

## 202. Pod в k8s убит OOMKilled — что делать?

1. `kubectl describe pod` — `Last State: Terminated, Reason: OOMKilled`.
2. Контейнер превысил **memory limit**.
3. **pprof heap** до рестарта (если успели снять) — что съело память: рост слайсов, кеш без eviction, утечка через global map.
4. Проверить **limits vs requests** — может limit слишком низкий для нормальной работы.
5. Go: GC не успевает при spike аллокаций; проверить `GOGC`, размер буферов, загрузку больших тел в память целиком.
6. Краткосрочно: поднять limit. Долгосрочно: профилировать heap, streaming вместо буферизации всего ответа.

## 203. Паника в горутине — что будет с сервисом?

Необработанная panic в горутине **убивает только эту горутину**, не весь процесс (в отличие от panic в main). Но: запрос оборвётся, ресурсы могут не освободиться, воркер исчезнет из пула.

**Решение:**

1. `defer recover()` в HTTP middleware и в worker goroutines — логировать stack, возвращать 500.
2. `net/http` уже ловит panic в handler.
3. В production: alert на panic в логах, stack trace в логере.
4. Не использовать recover для обычных ошибок — только для truly unexpected.

## 204. Connection pool к PostgreSQL исчерпан — почему и как починить?

Симптом: `pq: sorry, too many clients already`, таймауты на `db.Query`.

**Причины:**

1. **Утечка соединений** — `rows`/`stmt` не закрыты (`defer rows.Close()`).
2. **Слишком много реплик** сервиса × `MaxOpenConns` > лимит Postgres.
3. **Долгие транзакции** держат conn.
4. **Медленные запросы** — conn заняты, пул исчерпан.

**Решение:** проверить `pg_stat_activity`, закрывать rows, уменьшить `MaxOpenConns` на инстанс, PgBouncer, оптимизировать запросы, таймауты через `context`.

## 205. Клиент дважды отправил оплату — как защититься?

**Идемпотентность.**

1. Клиент шлёт `Idempotency-Key: uuid` в заголовке.
2. Сервер: таблица `idempotency_keys(key, response, status)` — при повторе вернуть сохранённый ответ.
3. В Kafka: `processed_events(event_id)` с unique constraint.
4. В БД: `UNIQUE(order_id)` + проверка статуса перед оплатой.

At-least-once доставка требует идемпотентной обработки на стороне consumer.

## 206. Два пользователя бронируют последнее место одновременно

Race condition без блокировки — оба прочитают `seats > 0`, оба запишут.

1. **Пессимистично**: `SELECT ... FOR UPDATE` в транзакции — второй ждёт, увидит 0 мест.
2. **Оптимистично**: `UPDATE ... WHERE count > 0` — проверить `RowsAffected == 1`, иначе retry.
3. **Redis**: `DECR` атомарно.

Выбор: высокая конкуренция на одну строку → `FOR UPDATE` или атомарный UPDATE; read-heavy → optimistic.

## 207. Thundering herd при протухании кеша

Тысяча запросов одновременно промахнулись в Redis — все пошли в БД.

1. **Jitter TTL** — разброс времени протухания (`TTL + random(0, 60s)`).
2. **Singleflight** (`golang.org/x/sync/singleflight`) — один запрос на ключ идёт в БД, остальные ждут.
3. **Probabilistic early refresh** — обновлять кеш до протухания с малой вероятностью.
4. **Mutex per key** в приложении.

В Go: `singleflight.Group` — стандартный паттерн.

## 208. Kafka consumer lag растёт — план действий

1. Lag по партициям — одна партиция (hot key) или все (мало consumers).
2. **Медленный handler** — pprof, время обработки одного сообщения, блокировки на БД.
3. **Rebalance** в логах — частые stop/start consumer group.
4. **Мало партиций** — увеличить партиции и инстансы consumer (до числа партиций).
5. **max.poll.interval.ms** — handler дольше интервала → consumer исключается из группы.
6. Batch processing + commit после batch.
7. Масштабировать consumers горизонтально.

## 209. Дубликаты сообщений из Kafka — как обрабатывать?

Kafka гарантирует at-least-once при retry. Consumer **должен быть идемпотентным**:

1. Хранить `event_id` в БД, `INSERT ... ON CONFLICT DO NOTHING`.
2. Проверка «уже обработано» перед side effect.
3. Сначала запись в БД (outbox/processed), потом side effect — или компенсация.
4. Для платежей — бизнес-ключ уникален в таблице.

Exactly-once сложно — проще идемпотентный handler.

## 210. 502/503 под нагрузкой — где искать

1. **Ingress/LB timeout** < время ответа сервиса — поднять timeout или ускорить сервис.
2. **Исчерпаны воркеры** — все goroutines заняты, очередь растёт.
3. **Upstream недоступен** — circuit breaker открыт, 503.
4. **k8s**: pod не ready, endpoints пустые.
5. **File descriptors** — лимит исчерпан.

Дальше: метрики RPS vs latency vs errors, трейс одного 502 запроса, load test.

## 211. Rolling update в k8s, но был даунтайм

**Причины:**

1. **Readiness probe** не проверяет реальную готовность (БД недоступна — pod в endpoints, но падает).
2. **`maxUnavailable: 100%`** или `replicas: 1`.
3. **PreStop hook** отсутствует — pod убит до drain из Service.
4. **Graceful shutdown** не реализован — SIGTERM, но процесс убит с активными запросами.

**Фикс:** readiness с проверкой зависимостей, `maxUnavailable: 0`, `maxSurge: 1`, `server.Shutdown`.

## 212. Сервис не подключается к БД после деплоя в k8s

1. **DNS** — hostname БД резолвится из pod (`nslookup`).
2. **Secret/ConfigMap** — правильные credentials, смонтированы в pod.
3. **NetworkPolicy** — egress к порту 5432 разрешён.
4. **Service name** — `postgres:5432` vs внешний host.
5. **SSL mode** — `sslmode=require` в prod.
6. Логи при старте — `db.Ping()` в readiness.

Локально работало через `localhost` — в k8s нужен cluster DNS.

## 213. N+1 запросов к БД — как обнаружить и исправить

Симптом: один HTTP-запрос → 100 SQL-запросов.

**Обнаружение:** query count per request, `pg_stat_statements`, трейс с N одинаковыми SELECT.

Типично: цикл `for _, user := range users { db.GetOrders(user.ID) }`.

**Фикс:**

1. `WHERE user_id IN (...)` одним запросом или JOIN.
2. Dataloader / batch load.
3. ORM: `Preload` в GORM.
4. В Go: один запрос + `map[userID][]Order` в памяти.

## 214. Как спроектировать rate limiter для публичного API?

Требования: лимит N req/min per API key / IP, 429 при превышении.

1. **In-memory** token bucket — просто, но не работает при N репликах без синхронизации.
2. **Redis** — `INCR` + `EXPIRE` или sliding window; общий лимит на кластер.
3. **API Gateway** (nginx, Kong) — лимит до приложения.

Middleware в Go: проверка до handler, заголовки `X-RateLimit-Remaining`. Учесть burst vs sustained rate.

## 215. Миграция с монолита на микросервисы без даунтайма (Strangler)

1. **Strangler Fig** — новый сервис через proxy/gateway, старый — остальное.
2. **Dual write** или **CDC (Debezium)** — синхронизация данных.
3. **Сравнение ответов** — shadow traffic, diff старого и нового API.
4. Постепенно переключать % трафика (feature flag, canary).
5. Откат всегда возможен.
6. Не делать big bang.

Для Go: отдельный сервис, общая БД на первом этапе допустима, потом выделение schema.

## 216. Медленный endpoint — пошаговая оптимизация

1. **Замерить** — benchmark, pprof, trace: CPU vs I/O bound.
2. **БД** — EXPLAIN, индексы, N+1, лишние поля.
3. **Кеш** — Redis для read-heavy, TTL, инвалидация.
4. **Аллокации** — `go test -bench -benchmem`, `sync.Pool` для буферов.
5. **Конкурентность** — параллельные вызовы через `errgroup`.
6. **Пагинация** вместо «вернуть всё».
7. Повторить замер — одно узкое место за раз, не гадать.
