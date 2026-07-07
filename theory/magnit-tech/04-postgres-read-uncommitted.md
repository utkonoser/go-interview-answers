# Magnit Tech: READ UNCOMMITTED в PostgreSQL

## M4. Почему в PostgreSQL нет первого уровня изоляции и как это решается?

### Формально уровень есть

```sql
BEGIN TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
```

PostgreSQL **принимает** синтаксис, но **игнорирует** его: внутри всегда работает как минимум **Read Committed**.

Проверка:

```sql
SHOW transaction_isolation;  -- read committed (по умолчанию)
```

### Почему нет настоящего READ UNCOMMITTED

PostgreSQL использует **MVCC** (Multi-Version Concurrency Control):

- Каждая транзакция видит **snapshot** зафиксированных версий строк
- Незакоммиченные изменения другой транзакции помечены «невидимыми» для читателей
- Читатель **не может** прочитать «грязную» версию строки — её просто нет в его snapshot

Dirty read противоречит модели MVCC: чтобы его разрешить, пришлось бы читать незакоммиченные tuple из чужой транзакции, ломая изоляцию и семантику отката.

### Что вместо этого

| Потребность | Решение в PG |
|-------------|--------------|
| Минимальные блокировки на чтение | MVCC: читатели не блокируют писателей |
| «Быстрое» чтение без ожидания lock | Read Committed + snapshot |
| Допустить грязное чтение | Не поддерживается намеренно |

### Read Committed в PG

- Каждый **отдельный statement** в транзакции видит snapshot на момент **начала этого statement**
- Отсюда non-repeatable read внутри одной транзакции — норма для RC

### Если нужна стабильная картина

```sql
BEGIN ISOLATION LEVEL REPEATABLE READ;
-- или
BEGIN ISOLATION LEVEL SERIALIZABLE;
```

### На собесе

Краткий ответ: «В PostgreSQL READ UNCOMMITTED = READ COMMITTED из-за MVCC; dirty read невозможен по дизайну». Длинный: объяснить snapshot, xmin/xmax, почему это плюс для ритейла (чеки, остатки) — не читаем незафиксированные списания.
