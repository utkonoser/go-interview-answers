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

## 2ГИС Application Review

Публичного репозитория с кодом секции **Application Review** нет — его выдают кандидату перед встречей ([как устроен найм](https://habr.com/ru/companies/2gis/articles/1034286/)). Задачи ниже воспроизводят темы из статьи: REST-сервис, счётчик запросов, бронирование + шина событий.

| # | Директория | Контекст | Что искать |
|---|------------|----------|------------|
| 11 | `11-request-metrics` | Счётчик RPS к ручке | Data race в middleware |
| 12 | `12-booking-outbox` | Бронь отеля + Kafka | Publish до commit, нет transactional outbox |

## Lamoda (склад / резерв)

Публичного кода **code review** с собеса нет. Зато известно классическое **тестовое на написание**: REST API склада — создать товар/склад, **reserve**, **release**, остатки ([lmd-tt](https://github.com/6jodeci/lmd-tt), [LamodaTestTask](https://github.com/zeroT4lant/LamodaTestTask), [lamoda-test-2023](https://github.com/adepte-myao/lamoda-test-2023)). На live-собесе — лайвкодинг, system design, сильная культура **360° code review** ([вакансии Lamoda](https://h.careers/job/98b46ce9-a942-40d4-bd08-e4d44ed4ddc3)).

| # | Директория | Контекст | Что искать |
|---|------------|----------|------------|
| 13 | `13-stock-reserve` | Резерв SKU на складе | TOCTOU reserve, Release без проверки, нет mutex |
| 14 | `14-stock-cache` | Кеш остатков (Redis) | Thundering herd, data race на map |
| 15 | `15-order-pricing` | Сумма заказа в checkout | `float64` для денег, потеря копеек |

Источники для вдохновения: [Avito playbook](https://github.com/avito-tech/playbook/blob/master/recruitment-and-office.md), [Habr: Code Review Horror Stories](https://habr.com/ru/articles/1031010/), [Habr: задачи WB/Ozon](https://habr.com/ru/articles/995600/), [2ГИС: собеседования Go](https://habr.com/ru/companies/2gis/articles/1034286/), [Lamoda: унификация Go-сервисов](https://habr.com/ru/companies/lamoda/articles/495344/).

В каждой папке два файла с build-тегами: по умолчанию собирается `code.go`, с `-tags=solution` — `solution.go`.
