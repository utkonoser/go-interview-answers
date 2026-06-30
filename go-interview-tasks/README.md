# Go Interview Tasks

Практические задачи по программированию на Go для подготовки к техническим собеседованиям. Каждая задача — в `strings/`, тесты — в `tests/`.

## Структура

```
go-interview-tasks/
├── strings/
│   ├── palindrome.go       # Проверка палиндрома
│   ├── reverse-string.go   # Разворот строки in-place
│   ├── backspace.go        # Сравнение строк с backspace
│   ├── remove-stars.go     # Удаление звёзд из строки
│   ├── dublicates.go       # Удаление соседних дубликатов
│   └── parentheses.go      # Проверка валидности скобок
└── tests/
    ├── palindrome_test.go
    ├── reverse-string_test.go
    ├── backspace_test.go
    ├── remove-stars_test.go
    ├── dublicates_test.go
    └── parentheses_test.go
```

## Задачи

| Файл | LeetCode | Описание |
|------|----------|----------|
| `palindrome.go` | #125 | Является ли строка палиндромом |
| `reverse-string.go` | #344 | Разворот строки на месте |
| `backspace.go` | #844 | Сравнение строк с учётом backspace |
| `remove-stars.go` | #2390 | Удаление `*` и символа слева |
| `dublicates.go` | #1047 | Удаление соседних одинаковых символов |
| `parentheses.go` | #20 | Валидность скобок `()[]{}` |

## Запуск тестов

```bash
# Из корня проекта (рекомендуется)
./scripts/check.sh

# Или вручную
cd go-interview-tasks
go test -v ./tests/...
go test -bench=. -benchmem ./tests/...
```

Отдельный пакет:

```bash
cd go-interview-tasks
go test -v ./tests/ -run TestPalindrome
```

## CI/CD

Тесты запускаются автоматически в GitHub Actions: [go-tests.yml](../.github/workflows/go-tests.yml)

[![Go Tests](https://github.com/utkonoser/interviews/workflows/Go%20Tests/badge.svg)](https://github.com/utkonoser/interviews/actions/workflows/go-tests.yml)

## Локальная проверка перед push

```bash
# Из корня проекта
./scripts/check.sh
```

Скрипт выполняет: `go test`, покрытие, бенчмарки, `go fmt`, `go vet`.
