# SQL-задачи

Проектирование схемы и написание запросов — типичный блок на собесах в Avito, Ozon, WB, Lamoda, Magnit, 2ГИС, Яндекс.

Диалект: **PostgreSQL** (на собесах чаще всего он; на livecoding иногда дают [sqlize.ru](https://sqlize.ru/) или песочницу в IDE).

## Как пользоваться

1. Открой `task.md` в директории задачи — условие и вопросы.
2. Изучи `schema.sql` и `seed.sql` (можно поднять локально: `psql -f schema.sql -f seed.sql`).
3. Напиши свои запросы / схему в отдельном файле или на бумаге.
4. Сверься с `solution.sql`.

```bash
# пример локальной проверки (нужен PostgreSQL)
createdb interview_sql
psql interview_sql -f sql/01-avito-deals-analytics/schema.sql
psql interview_sql -f sql/01-avito-deals-analytics/seed.sql
psql interview_sql -f sql/01-avito-deals-analytics/solution.sql
```

## Задачи

| # | Директория | Компания / домен | Фокус |
|---|------------|------------------|-------|
| 1 | `01-avito-deals-analytics` | Avito / Ozon (маркетплейс объявлений) | Аналитика сделок, воронка просмотров |
| 2 | `02-wb-stock-replenishment` | WB / Lamoda (склад, заказы) | Проектирование остатков + запрос на дефицит |
| 3 | `03-2gis-bookings` | 2ГИС (бронирование отелей) | Пересечение дат, загрузка номеров |
| 4 | `04-magnit-loyalty` | Магнит (ритейл, лояльность) | Журнал баллов + RFM-подобная аналитика |

На собесе обычно 20–40 минут: 1–2 запроса + обсуждение индексов, изоляции, граничных случаев.
