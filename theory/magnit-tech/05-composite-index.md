# Magnit Tech: составной индекс

## M5. Как выбирать порядок колонок в составном индексе?

Индекс `(a, b, c)` — B-tree сначала сортирует по `a`, внутри одинаковых `a` — по `b`, и т.д.

### Правило левого префикса

Индекс `(city, status, created_at)` используется для:

- `WHERE city = ?` — да
- `WHERE city = ? AND status = ?` — да
- `WHERE city = ? AND status = ? ORDER BY created_at` — да, часто без sort
- `WHERE status = ?` — **нет** (пропущен `city` слева)
- `WHERE city = ? AND created_at > ?` — `city` да, `created_at` как range — дальше по индексу хуже

### Порядок колонок

1. **Равенство (=, IN)** — слева, в порядке селективности (часто: более селективный первым, если оба в WHERE)
2. **Диапазон (>, <, BETWEEN)** — после equality-колонок; обычно **одна** range-колонка в индексе эффективна
3. **ORDER BY** — колонки сортировки в конце, если совпадают с направлением индекса

### Пример: лента товаров магазина

```sql
-- Запрос
SELECT * FROM products
WHERE store_id = 42 AND category = 'dairy' AND price BETWEEN 100 AND 500
ORDER BY updated_at DESC
LIMIT 20;

-- Индекс
CREATE INDEX idx_products_store_cat_price_updated
    ON products (store_id, category, price, updated_at DESC);
```

`store_id`, `category` — equality; `price` — range; `updated_at` — sort (если planner использует index scan + filter).

### Covering index (INCLUDE)

```sql
CREATE INDEX idx_orders_user_created
    ON orders (user_id, created_at DESC)
    INCLUDE (total_kopecks, status);
```

Index-only scan, если все колонки SELECT в индексе.

### Антипаттерны

- `(created_at, store_id)` при фильтре только по `store_id` — индекс бесполезен
- Дублировать индексы `(a,b)` и `(a,b,c)` — иногда `(a,b,c)` покрывает `(a,b)`
- Ставить low-cardinality поле первым (`gender`, `is_active`) — индекс слабый

### EXPLAIN

Всегда проверять: `EXPLAIN (ANALYZE, BUFFERS) SELECT ...` — Index Scan vs Seq Scan, Sort отдельно или нет.
