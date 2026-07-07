# Avito: code review — указатели и горутины в цикле

## A4. Указатели: что выведет программа?

```go
type Person struct { Name string }

func changeName(person *Person) {
    person = &Person{Name: "Alice"} // меняем копию указателя
}

func main() {
    person := &Person{Name: "Bob"}
    fmt.Println(person.Name) // Bob
    changeName(person)
    fmt.Println(person.Name) // Bob
}
```

**Почему:** в Go **нет pass-by-reference** для параметров. `changeName` получает **копию** указателя; присваивание `person = &Person{...}` меняет только локальную копию.

**Исправление — менять поле по указателю:**

```go
func changeName(person *Person) {
    person.Name = "Alice"
}
// или *person = Person{Name: "Alice"}
```

---

## A4b. Горутины в цикле: максимальное чётное

```go
var max int
for i := 1000; i > 0; i-- {
    go func() {
        if i%2 == 0 && i > max {
            max = i
        }
    }()
}
fmt.Printf("Maximum is %d", max)
```

**Три проблемы (ожидают на собесе):**

### 1. Loop variable capture
Все goroutines видят **одно** `i` (Go ≤1.21) или последнее значение. `go vet`: `loop variable i captured by func literal`.

**Fix:** `go func(v int) { ... }(i)`

### 2. Нет синхронизации завершения
`main` может выйти до старта goroutines.

**Fix:** `sync.WaitGroup`

### 3. Data race на `max`
Read-modify-write не атомарен. Atomic только на `max = i` **не спасает** — нужно лочить **и сравнение, и запись** (check-then-act).

**Fix:**

```go
mu.Lock()
if v%2 == 0 && v > max {
    max = v
}
mu.Unlock()
```

**Инструменты:** `go vet`, `go test -race`, `-race` в CI.

**Follow-up:** зачем `defer wg.Done()` vs прямой вызов; low-cost defers (Go 1.14+); можно ли читать `max` без lock — нет, если параллельно пишут.
