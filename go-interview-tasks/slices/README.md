# Слайсы: «что выведет код?»

Классические задачи с собесов (Avito, Ozon, Яндекс, WB): общий backing array, `append`, `copy`, передача в функцию.

## Как заниматься

1. Открой `puzzles.go`, найди `Puzzle01`, `Puzzle02`, …
2. Не запуская код, запиши ответ: содержимое слайса, `len`, `cap`.
3. Проверь себя:

```bash
cd go-interview-tasks
go test -v ./slices/...
```

Тесты падают, если ты ошибся — в выводе будет `want` vs `got`.

## Список головоломок

| # | Функция | Суть |
|---|---------|------|
| 1 | `Puzzle01AppendInPlace` | `append` в subslice при запасе `cap` — меняет родителя |
| 2 | `Puzzle02AppendRealloc` | `append` за пределы `cap` — новый массив, родитель цел |
| 3 | `Puzzle03AssignShare` | `b := a` — общий backing array |
| 4 | `Puzzle04CopyPartial` | `copy` копирует `min(len(dst), len(src))` элементов |
| 5 | `Puzzle05CopyOverlap` | `copy(s[1:], s[2:])` — сдвиг внутри одного массива |
| 6 | `Puzzle06FullSliceExpr` | `s[low:high:max]` ограничивает `cap`, `append` не трогает родителя |
| 7 | `Puzzle07AppendNoAssign` | `append` в функции без присваивания — длина снаружи не растёт |
| 8 | `Puzzle08MutateIndex` | Индекс в функции меняет общий массив |

Разбор — в комментариях к тестам в `puzzles_test.go`.
