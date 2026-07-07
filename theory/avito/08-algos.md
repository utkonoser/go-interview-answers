# Avito: алгоритмические задачи

## A8. Чемпионат по шагам

**Условие:** по дням соревнований — списки `{userId, steps}`. Найти userId с **максимальной суммой шагов**, участвовавших **во все дни**. Несколько победителей возможно; если никто не был все дни — пустой результат.

**Сложность:** O(N) по всем записям или O(N×K) по дням × участников дня.

### Решение

```go
type stat struct {
    days   int
    steps  int
}

func getChampions(statistics [][]Statistic) Result {
    if len(statistics) == 0 {
        return Result{}
    }
    totalDays := len(statistics)
    m := make(map[int]stat, len(statistics[0]))

    for _, day := range statistics[0] {
        m[day.UserID] = stat{days: 1, steps: day.Steps}
    }
    for d := 1; d < totalDays; d++ {
        for _, rec := range statistics[d] {
            if s, ok := m[rec.UserID]; ok {
                s.days++
                s.steps += rec.Steps
                m[rec.UserID] = s
            }
        }
    }

    maxSteps := -1
    for _, s := range m {
        if s.days == totalDays && s.steps > maxSteps {
            maxSteps = s.steps
        }
    }
    if maxSteps < 0 {
        return Result{}
    }

    var ids []int
    for id, s := range m {
        if s.days == totalDays && s.steps == maxSteps {
            ids = append(ids, id)
        }
    }
    sort.Ints(ids)
    return Result{UserIDs: ids, Steps: maxSteps}
}
```

**Corner cases:**
- несколько чемпионов с одной суммой;
- нет участника все дни → `{}`;
- не класть в map участников только из поздних дней (оптимизация: seed из day 0).

**Anti-patterns:** `includes` на каждом шаге O(K²); сортировка всего массива для ответа.

---

## A8b. Суммарная неудовлетворённость покупателей

**Условие:** `goods[]`, `buyerNeeds[]`. Покупатель берёт товар с **минимальной** `|good - need|` (товаров бесконечно). Сумма разниц по всем покупателям.

**Пример:** `goods=[8,3,5]`, `needs=[5,6]` → `0 + 1 = 1`.

### Варианты

| Подход | Сложность |
|--------|-----------|
| Полный перебор | O(n×m) |
| Sort + two pointers | O(n log n + m log m) |
| Sort goods + binary search per need | O(n log n + m log n) |

### Binary search (типичное ожидание)

```go
func findTotalDissatisfaction(goods, needs []uint) uint {
    sort.Slice(goods, func(i, j int) bool { return goods[i] < goods[j] })
    var total uint
    for _, need := range needs {
        total += abs(closest(goods, need), need)
    }
    return total
}

func closest(sorted []uint, target uint) uint {
    lo, hi := 0, len(sorted)-1
    for lo+1 < hi {
        mid := (lo + hi) / 2
        if sorted[mid] < target {
            lo = mid
        } else {
            hi = mid
        }
    }
    if abs(sorted[lo], target) <= abs(sorted[hi], target) {
        return sorted[lo]
    }
    return sorted[hi]
}
```

**Two pointers** работает если **оба** массива отсортировать и двигать указатель по goods для каждого need monotonically.

**Ошибки:** забыть tie-breaking при равной разнице; overflow (в задаче uint); не отсортировать goods перед бинпоиском.
