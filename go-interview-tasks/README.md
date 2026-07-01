# Практические задачи

Задачи на Go для тренировки перед собесами.

- **Алгоритмы** — `strings/`, тесты в `tests/`
- **Code review** — `code-review/` (01–07 базовые, 08–10 комплексные, 11–12 2ГИС, 13–15 Lamoda)

```bash
./scripts/check.sh
```

Только code review (решения):

```bash
cd go-interview-tasks
go test -tags=solution -race ./code-review/tests/...
```
