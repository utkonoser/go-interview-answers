# Практические задачи

Задачи на Go для тренировки перед собесами.

- **Алгоритмы** — `strings/`, тесты в `tests/`
- **Code review** — `code-review/` (базовые 01–07 + комплексные 08–10, см. README)

```bash
./scripts/check.sh
```

Только code review (решения):

```bash
cd go-interview-tasks
go test -tags=solution -race ./code-review/tests/...
```
