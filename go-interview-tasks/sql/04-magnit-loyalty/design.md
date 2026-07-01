# Эталон проектирования (часть 1)

## Журнал + баланс

`point_ledger` — source of truth для истории; `loyalty_cards.points_balance` — кеш для быстрого чтения на кассе.

В одной транзакции:

```sql
INSERT INTO point_ledger (customer_id, purchase_id, delta, reason)
VALUES ($customer, $purchase, $accrual, 'accrual_1pct');

UPDATE loyalty_cards
SET points_balance = points_balance + $accrual
WHERE customer_id = $customer;
```

## Начисление 1%

`accrual := purchase.total_kopecks / 100` (целочисленное деление).

## Списание без двойного списания

```sql
-- UNIQUE (purchase_id, reason) WHERE reason = 'redeem' — одна redemption на чек

INSERT INTO point_ledger (customer_id, purchase_id, delta, reason)
VALUES ($c, $p, -$redeem, 'redeem');

UPDATE loyalty_cards
SET points_balance = points_balance - $redeem
WHERE customer_id = $c AND points_balance >= $redeem;
-- проверить row count; иначе rollback
```

Идемпотентность начисления:

```sql
UNIQUE (purchase_id, reason)  -- reason = 'accrual_1pct'
```

Повторная обработка чека → unique violation → no-op.

## Верификация баланса

Периодический job: `points_balance = SUM(delta)`; расхождения в алерт.

На собесе плюс: event sourcing только через ledger без кеша — медленнее на чтение в POS.
