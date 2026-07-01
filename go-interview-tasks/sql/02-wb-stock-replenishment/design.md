# Эталон проектирования (часть 1)

## Таблица `warehouse_stock`

Уже в `schema.sql`:

| Колонка | Смысл |
|---------|--------|
| `quantity` | физический остаток на полке |
| `reserved` | зарезервировано под `pending`/`confirmed` заказы |

Инвариант: `0 <= reserved <= quantity`.  
Свободно: `quantity - reserved`.

## Reserve (в транзакции)

```sql
BEGIN;

SELECT quantity, reserved
FROM warehouse_stock
WHERE warehouse_id = $1 AND product_id = $2
FOR UPDATE;

-- в приложении: if quantity - reserved < $qty → rollback, insufficient stock

UPDATE warehouse_stock
SET reserved = reserved + $qty
WHERE warehouse_id = $1 AND product_id = $2
  AND quantity - reserved >= $qty;

-- affected rows = 0 → не хватило остатка

INSERT INTO order_items (...) VALUES (...);

COMMIT;
```

## Release (отмена или отгрузка)

- **Отмена заказа:** `reserved -= qty` по всем позициям, `orders.status = 'cancelled'`.
- **Отгрузка:** `quantity -= qty`, `reserved -= qty`, статус `shipped`.

Оба шага — в одной транзакции с блокировкой строки `warehouse_stock`.

## Почему не только `SUM(order_items)`?

- Снимок остатков нужен для O(1) проверки на кассе/в API; агрегат по всем заказам дорог и гоняется с историей.
- Резерв в `warehouse_stock` — источник правды для **доступно к продаже**; заказы — журнал намерений.

## Индексы

- `warehouse_stock (product_id)` — поиск по SKU через join.
- `orders (warehouse_id, status, created_at)` — отчёты.
- `order_items (product_id)` — спрос по товару.

Опционально: `stock_movements` (audit log) для разборов инцидентов — на собесе плюс, не обязательно.
