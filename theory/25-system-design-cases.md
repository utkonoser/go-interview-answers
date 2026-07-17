# Системный дизайн — топ-5 кейсов

Самые частые задачи на секции SD (backend, middle/senior). Подходят для любой компании — ритейл, финтех, соцсети.

На собесе ждут **структуру мышления**, не «правильную» диаграмму: уточнить требования → оценки → API → данные → компоненты → узкие места → trade-offs.

## Как отвечать (скелет на 35–45 мин)

1. **Уточнить scope** — MVP vs v2 (5 мин).
2. **Оценки** — users, RPS, storage, latency SLA (5 мин).
3. **API** — 2–4 метода, sync vs async (5 мин).
4. **Хранение и кэш** (10 мин).
5. **Архитектура и поток данных** (10 мин).
6. **Масштабирование, отказы, мониторинг** (5 мин).
7. **Что упростил и почему** (2 мин).

---

## 210. URL Shortener (сокращатель ссылок)

**Постановка:** `POST /shorten` → короткая ссылка; `GET /{code}` → 302 на оригинал. Как bit.ly.

**Уточнить:** кастомные алиасы? срок жизни ссылки? аналитика кликов? нужна ли авторизация?

**Оценки:** 100M новых URL/мес, read:write ≈ 100:1, redirect p99 < 50 ms.

**API:**
- `POST /api/v1/urls` — `{ "long_url": "..." }` → `{ "short_url": "https://x.co/abc12" }`
- `GET /{code}` — 302 Redirect

**Генерация кода:**
- Base62 от auto-increment id (просто, предсказуемо — можно добавить salt).
- Или hash(long_url) — коллизии → retry.
- Или случайная строка 6–8 символов + unique index.

**Хранение:**
```
urls(code PK, long_url, user_id, created_at, expires_at)
```
- Шардирование по `hash(code)` при росте.
- **Redis** — hot codes для redirect (cache-aside).

**Redirect path:**
```
Client → CDN/LB → Redirect Service → Redis → (miss) → DB → 302
```

**Аналитика кликов (async):**
- Событие в Kafka → агрегаты в ClickHouse/PostgreSQL.
- Не блокировать redirect.

**Узкие места:** hot URL (вирусная ссылка) — CDN кэширует 302; счётчик кликов — async, не в hot path.

---

## 211. Rate Limiter (ограничение RPS)

**Постановка:** не больше N запросов в минуту на user/IP/API-key. 429 при превышении.

**Уточнить:** лимит глобальный или per-user? sliding window или fixed? распределённый кластер?

**Алгоритмы:**

| Алгоритм | Плюсы | Минусы |
|----------|-------|--------|
| **Token bucket** | burst допустим | чуть сложнее |
| **Fixed window** | просто | spike на границе окна |
| **Sliding window log** | точно | память |
| **Sliding window counter** | баланс | — |

**Распределённый вариант (Redis):**
```
INCR rate:{user_id}:{window}
EXPIRE ... TTL = window_size
if count > limit → 429
```
Или Lua-скрипт для атомарности.

**Где ставить:** API Gateway / sidecar / middleware в Go (`http.Handler` chain).

**Trade-offs:** Redis single point → Redis Cluster; при падении Redis — fail-open (пропускать) vs fail-closed (блокировать) — озвучить выбор.

**Подробнее:** `24-scenarios.md` (207), `magnit-tech/06-rate-limiter.md`.

---

## 212. Notification System (push / email / SMS)

**Постановка:** сервис рассылает уведомления пользователям. Каналы: push, email, SMS. Приоритеты, retry, не слать дубли.

**Уточнить:** real-time или допустима задержка? объём (1M push/час)? пользователь выбирает каналы?

**Архитектура:**
```
API → Notification Service → Queue (Kafka/SQS)
                           → Workers per channel (push, email, sms)
                           → External providers (FCM, SendGrid, Twilio)
```

**Модель данных:**
```
notifications(id, user_id, template_id, payload, status, created_at)
user_preferences(user_id, channel, enabled)
```

**Ключевые паттерны:**
- **Очередь** — decouple API от медленных провайдеров.
- **Idempotency key** — не отправить дважды при retry.
- **Template service** — рендер тела по шаблону.
- **Rate limit per provider** — SMS дорогой, throttle отдельно.
- **DLQ** — после N retry в dead letter queue + алерт.

**Fan-out:** одно событие «новый заказ» → 1M подписчиков — не в одном запросе; batch consumers, partition по user_id.

**Мониторинг:** delivery rate, latency per channel, DLQ depth.

---

## 213. News Feed / лента (Twitter, Instagram)

**Постановка:** пользователь публикует пост; подписчики видят ленту. `POST /tweet`, `GET /feed`.

**Уточнить:** лента хронологическая или ranked (ML)? celebrity с 10M подписчиков?

**Два подхода:**

| | Fan-out on write | Fan-out on read |
|--|------------------|-----------------|
| Идея | При посте пишем в ленту каждого подписчика | При чтении собираем посты подписок |
| Плюсы | быстрый read | дешёвый write |
| Минусы | дорогой write для celebrity | медленный read |

**Гибрид (как в проде):**
- Обычные юзеры — **fan-out on write** → Redis `feed:{user_id}` (sorted set by timestamp).
- Celebrity — **fan-out on read** + кэш их постов отдельно.

**Хранение:**
- PostgreSQL / Cassandra — посты, подписки (follows).
- Redis — precomputed feeds.
- CDN — медиа.

**Оценки:** 300M DAU, 500M постов/день, read >> write.

**Пагинация:** cursor-based (`since_id` / `max_id`), не offset на больших лентах.

---

## 214. Chat / Messenger (WhatsApp, Telegram-lite)

**Постановка:** 1:1 и групповые чаты, доставка в реальном времени, история сообщений, online status.

**Уточнить:** только online или offline тоже? end-to-end encryption? размер группы?

**Компоненты:**
- **WebSocket / long polling** — real-time доставка.
- **Chat Service** — логика чатов, membership.
- **Message Service** — persist + порядок сообщений.
- **Presence Service** — online/offline (Redis TTL heartbeat).

**Доставка online:**
```
Sender → WS Gateway → Message Service → DB
                    → WS Gateway получателя (если online)
                    → Push Service (если offline)
```

**Хранение сообщений:**
- Шард по `chat_id` или `user_id`.
- Пагинация: `GET /chats/{id}/messages?before={msg_id}&limit=50`.

**Порядок:** `server_timestamp` + `message_id` (snowflake); при конфликте — клиент merge.

**Групповой чат:** не пушить 10k WS — fan-out через message queue; для больших групп — только pull.

**Масштабирование WS:** sticky sessions по user_id или shared pub/sub (Redis) между gateway-нодами.

---

## Шпаргалка (встречается во всех кейсах)

| Тема | Что сказать |
|------|-------------|
| Масштабирование | Stateless + LB + horizontal scale |
| Кэш | Cache-aside, TTL, инвалидация |
| БД | Read replicas, шардирование по ключу |
| Очередь | Async, retry, DLQ, at-least-once + idempotency |
| CAP | Озвучить выбор: CP (деньги) vs AP (лента, аналитика) |
| Мониторинг | RED / USE, алерты на lag и error rate |

## Связанные файлы

- Паттерны, circuit breaker → `11-system-design.md`
- Troubleshooting → `24-scenarios.md`
- Kafka → `20-kafka.md`
- Кэш, thundering herd → `24-scenarios.md` (200)
- k8s, деплой → `21-docker-kubernetes.md`
