# Magnit Tech: panic и error

## M2. Как работают panic в Go и чем отличаются от error?

### Error — обычное значение

```go
func readConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }
    // ...
}
```

- Возвращается явно, обрабатывается `if err != nil`
- Не раскручивает стек
- Ожидаемые сбои: файл не найден, таймаут, валидация

### Panic — аварийная остановка горутины

```go
panic("invariant violated")
```

- Раскручивает стек, выполняет **defer** по пути
- Если нет `recover` — **горутина падает** (main упадёт, HTTP handler в net/http ловится middleware)
- Для «это невозможно» / баг программиста, не для бизнес-ошибок

### recover

```go
func handler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if v := recover(); v != nil {
            log.Printf("panic: %v\n%s", v, debug.Stack())
            http.Error(w, "internal error", 500)
        }
    }()
    // ...
}
```

`recover()` работает **только в defer** той же горутины, где был panic.

### Ключевые отличия

| | `error` | `panic` |
|---|---------|---------|
| Когда | Ожидаемый сбой | Баг / невосстановимо в месте вызова |
| Контроль | Вызывающий решает | Раскрутка стека или recover |
| Горутина | Продолжает работу | Умирает без recover |
| В HTTP API | 4xx/5xx в ответе | 500 + лог (если поймали) |

### Panic в горутине

Panic в **дочерней** горутине **не убивает** main и другие горутины — только её. Сервис «жив», но запрос/воркер оборван, возможна утечка ресурсов. В production: recover в middleware и worker pool.

### Когда что на собесе

- Валидация входа, БД недоступна, 404 → **error**
- `nil` pointer по логике «не бывает», index out of range после бага → можно panic в dev, в prod лучше error + метрика
- **Не** использовать panic для flow control
