# Observability (метрики, логи, трейсы)

## 194. Что такое observability и из чего она состоит?

Способность понять внутреннее состояние системы по внешним сигналам. Три столпа: **метрики** (числа во времени), **логи** (события с контекстом), **трейсы** (путь запроса через сервисы). Плюс **профилирование** (pprof) для глубокого анализа. Цель: быстро отвечать «почему медленно / почему упало» в production.

## 195. Чем метрики отличаются от логов?

**Метрики** — агрегированные числа: RPS, latency p99, error rate, CPU. Дешёвые в хранении, хороши для алертов и дашбордов. **Логи** — дискретные события: «заказ 123 создан», stack trace. Дороже в объёме, нужны для детальной отладки. Правило: алертить по метрикам, расследовать по логам и трейсам.

## 196. Что такое RED и USE методологии?

**RED** (для сервисов): **R**ate (запросов/сек), **E**rrors (доля ошибок), **D**uration (латентность). **USE** (для ресурсов): **U**tilization (% CPU/диска), **S**aturation (очередь, ожидание), **E**rrors. Минимальный набор метрик для любого Go HTTP/gRPC сервиса: `http_requests_total`, `http_request_duration_seconds`, error counter по status code.

## 197. Как устроен Prometheus и pull-модель?

Prometheus **скрейпит** (pull) `/metrics` endpoint каждые N секунд. Метрики в формате text exposition. Хранит time series в TSDB. **Alertmanager** — алерты по правилам. В Go: `prometheus/client_golang` — Counter, Gauge, Histogram, Summary. Histogram для latency с buckets. Grafana — визуализация. Push — только для batch jobs (Pushgateway).

## 198. Какие метрики собирать у Go-сервиса?

Обязательно: RPS, latency (histogram p50/p95/p99), error rate по endpoint. Runtime: `go_goroutines`, `go_memstats_*`, GC pause — из `prometheus/client_golang` collectors. Бизнес-метрики: заказы/сек, размер очереди. Не раздувать cardinality: не label'ить `user_id` на каждый запрос — взрыв time series.

## 199. Что такое distributed tracing?

Отслеживание одного запроса через цепочку сервисов. **Trace** — весь путь, **span** — один участок (HTTP call, DB query). **Trace ID** прокидывается через заголовки (`traceparent`, Jaeger). В Go: OpenTelemetry SDK → Jaeger/Tempo. Показывает: где потратили 2 из 3 секунд, какой сервис тормозит.

## 200. Структурированное логирование в Go

Текстовые логи плохо парсятся. **Структурированные** — JSON/key-value: `slog.Info("order created", "order_id", id, "user_id", uid)`. Поля фильтруются в Loki/ELK. Стандарт: `log/slog` (Go 1.21+). Обязательно: **request_id** / trace_id в каждой строке. Уровни: Debug, Info, Warn, Error. Не логировать секреты и PII.

## 201. Как профилировать Go-сервис в production?

`net/http/pprof` на admin-порту (не публичном): `/debug/pprof/profile` (CPU), `/debug/pprof/heap`. Снять: `go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30`. Continuous profiling: Pyroscope, Grafana Phlare. Метрики + трейсы находят «что медленно», pprof — «почему в коде». Включать с низким overhead, защищать auth/network policy.
