# Практические задачи

Задачи на Go для тренировки перед собесами.

- **Алгоритмы** — `strings/`, тесты в `tests/`
- **Слайсы** — `slices/` (головоломки «что выведет?», append/copy/aliasing)
- **Live-coding** — `live-coding/` (LRU, worker pool, rate limiter, graceful shutdown)
- **Code review** — `code-review/` (01–07 базовые, 08–10 комплексные, 11–12 2ГИС, 13–15 Lamoda)
- **SQL** — `sql/` (Avito/Ozon, WB/Lamoda, 2ГИС, Магнит)

```bash
./scripts/check.sh
```

Только code review (решения):

```bash
cd go-interview-tasks
go test -tags=solution -race ./code-review/tests/...
```

Live-coding:

```bash
go test ./live-coding/tests/...
```

Слайсы (самопроверка):

```bash
go test -v ./slices/...
```
