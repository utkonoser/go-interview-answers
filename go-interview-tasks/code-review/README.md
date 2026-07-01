# Code review

Задачи на разбор чужого Go-кода: найти баги, уязвимости и антипаттерны.

## Как пользоваться

1. Открой `code.go` в директории задачи — это код «коллеги» на ревью.
2. Запиши замечания: баги, race, утечки, ошибки API, читаемость.
3. Сверься с `solution.go` — эталон после ревью.
4. Прогони тесты исправленной версии:

```bash
cd go-interview-tasks
go test -tags=solution ./code-review/tests/...
```

Для поиска data race:

```bash
go test -tags=solution -race ./code-review/tests/...
```

## Задачи

| # | Директория | Тема |
|---|------------|------|
| 1 | `01-race-counter` | Data race на счётчике |
| 2 | `02-goroutine-leak` | Утечка горутины, канал не закрыт |
| 3 | `03-waitgroup-misuse` | `WaitGroup.Add` внутри горутины |
| 4 | `04-http-client` | HTTP-клиент без таймаута |
| 5 | `05-concurrent-map` | Конкурентная запись в map |
| 6 | `06-defer-in-loop` | `defer` внутри цикла |
| 7 | `07-error-handling` | Неправильный порядок `defer` и проверки ошибки |

## Комплексные (Middle+, стиль Avito/Ozon/WB)

| # | Директория | Контекст | Что искать |
|---|------------|----------|------------|
| 8 | `08-ad-click-tracker` | Ad-tech, дедуп кликов | Mutex, TOCTOU, канал без consumer, value copy |
| 9 | `09-wallet-transfer` | Баланс (паттерн Avito) | Race на transfer, атомарность debit |
| 10 | `10-parallel-search` | Параллельный поиск (WB) | Захват цикла, гонка на slice/int |

Источники для вдохновения: [Avito playbook](https://github.com/avito-tech/playbook/blob/master/recruitment-and-office.md) (секция code review на техсобесе), [Habr: Code Review Horror Stories](https://habr.com/ru/articles/1031010/) (~150 строк сервиса, 20+ багов), [Habr: задачи WB/Ozon](https://habr.com/ru/articles/995600/).

В каждой папке два файла с build-тегами: по умолчанию собирается `code.go`, с `-tags=solution` — `solution.go`.
