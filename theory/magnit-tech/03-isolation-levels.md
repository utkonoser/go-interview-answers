# Magnit Tech: уровни изоляции

## M3. Примеры аномалий на разных уровнях изоляции (SQL)

Таблица для демо:

```sql
CREATE TABLE accounts (id INT PRIMARY KEY, balance INT);
INSERT INTO accounts VALUES (1, 1000);
```

Две сессии: **T1** и **T2**. Запускать в двух окнах `psql`.

### Dirty Read — чтение незафиксированных данных

**Уровень:** Read Uncommitted (в PostgreSQL фактически не даёт dirty read — см. M4).

```sql
-- T1
BEGIN;
UPDATE accounts SET balance = 500 WHERE id = 1;  -- не COMMIT

-- T2 (Read Uncommitted в MySQL)
SELECT balance FROM accounts WHERE id = 1;  -- может увидеть 500

-- T1
ROLLBACK;  -- T2 читал «грязные» 500, которых не было
```

### Non-Repeatable Read

**Уровень:** Read Committed допускает.

```sql
-- T1
BEGIN;  -- PostgreSQL default: READ COMMITTED
SELECT balance FROM accounts WHERE id = 1;  -- 1000

-- T2
BEGIN;
UPDATE accounts SET balance = 800 WHERE id = 1;
COMMIT;

-- T1 — повторное чтение
SELECT balance FROM accounts WHERE id = 1;  -- 800 (было 1000)
COMMIT;
```

Один и тот же SELECT в транзакции T1 дал разный результат.

### Phantom Read

**Уровень:** Repeatable Read в PostgreSQL частично закрывает, Serializable — полностью.

```sql
CREATE TABLE orders (id SERIAL, account_id INT, amount INT);
INSERT INTO orders (account_id, amount) VALUES (1, 100);

-- T1
BEGIN ISOLATION LEVEL REPEATABLE READ;
SELECT SUM(amount) FROM orders WHERE account_id = 1;  -- 100

-- T2
INSERT INTO orders (account_id, amount) VALUES (1, 200);
COMMIT;

-- T1
SELECT SUM(amount) FROM orders WHERE account_id = 1;
-- в PG Repeatable Read: snapshot — всё ещё 100
-- в Serializable / при другом движке — может быть 300 (phantom)
COMMIT;
```

### Сводка

| Аномалия | Read Uncommitted | Read Committed | Repeatable Read | Serializable |
|----------|------------------|----------------|-----------------|--------------|
| Dirty read | ✓ | ✗ | ✗ | ✗ |
| Non-repeatable read | ✓ | ✓ | ✗ | ✗ |
| Phantom read | ✓ | ✓ | ✓* | ✗ |

\* PostgreSQL RR блокирует phantom за счёт snapshot isolation.

### Go: уровень в транзакции

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})
```

Для перевода денег / резерва SKU на складе Magnit — минимум **Repeatable Read** или `SELECT FOR UPDATE`.
