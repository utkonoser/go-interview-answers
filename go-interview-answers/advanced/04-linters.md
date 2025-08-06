# 17. Какой у вас любимый линтер?

## Ответ

**golangci-lint** — самый популярный и мощный линтер для Go.

### Почему golangci-lint?

- Объединяет множество линтеров в одном инструменте
- Настраиваемые правила
- Интеграция с CI/CD
- Высокая производительность

### Основные линтеры в составе:

```go
// golint - проверяет стиль кода
func BadFunctionName() {} // ❌

// govet - статический анализ
func example() {
    var x int
    fmt.Printf("%s", x) // ❌ неправильный формат
}

// errcheck - проверяет обработку ошибок
func bad() {
    os.Open("file.txt") // ❌ ошибка не обработана
}

// staticcheck - современный анализатор
func unused() {
    var x int // ❌ неиспользуемая переменная
}
```

### Конфигурация .golangci.yml:

```yaml
linters:
  enable:
    - golint
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell
    - gosec

run:
  timeout: 5m
  tests: true

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

### Интеграция с CI:

```yaml
# .github/workflows/lint.yml
name: Lint
on: [push, pull_request]
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```

### Альтернативы:

- **gofmt** - форматирование кода
- **go vet** - встроенный анализатор
- **staticcheck** - современный статический анализатор
- **revive** - быстрая замена golint

### Лучшие практики:

1. Настройте pre-commit хуки
2. Интегрируйте в CI/CD pipeline
3. Используйте разные уровни строгости для dev/prod
4. Регулярно обновляйте правила 