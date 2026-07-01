# Эталон проектирования (часть 1)

## Инвариант

Для одного `room_id` интервалы `[check_in, check_out)` у броней со статусом `confirmed` **не пересекаются**.

## PostgreSQL: EXCLUDE

```sql
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE bookings
    ADD CONSTRAINT bookings_no_overlap
    EXCLUDE USING gist (
        room_id WITH =,
        daterange(check_in, check_out, '[)') WITH &&
    )
    WHERE (status = 'confirmed');
```

`[)` — check_in входит, check_out — день выезда, не ночь.

## Отмена

`status = 'cancelled'` — строка выпадает из partial constraint (`WHERE confirmed`). Слот освобождается без удаления истории.

## Альтернатива: `room_nights`

Таблица `(room_id, night DATE)` с `UNIQUE (room_id, night)` — проще для отчётов по дням, тяжелее при длинных бронях (много строк). На собесе достаточно обсудить trade-off.

## Индексы

- `(room_id, check_in)` WHERE `status = 'confirmed'` — поиск конфликтов при вставке.
- `(place_id)` на `rooms` через join к отелям — отчёты по городу.
