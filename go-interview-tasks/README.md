# 💻 Go Interview Tasks

Практические задачи по программированию на Go для подготовки к техническим собеседованиям. Каждая задача содержит оптимизированное решение и comprehensive тесты.

## 📁 Структура

```
go-interview-tasks/
├── strings/                    # Задачи по работе со строками
│   ├── palindrome.go          # Проверка палиндрома
│   ├── reverse-string.go      # Разворот строки in-place
│   ├── backspace.go           # Сравнение строк с backspace
│   ├── remove-stars.go        # Удаление звезд из строки
│   ├── remove-duplicates.go   # Удаление дубликатов
│   └── parentheses.go         # Проверка валидности скобок
├── tests/                     # Тесты для всех задач
│   ├── palindrome_test.go     # Тесты палиндрома
│   ├── reverse-string_test.go # Тесты разворота строки
│   ├── backspace_test.go      # Тесты backspace
│   ├── remove-stars_test.go   # Тесты удаления звезд
│   ├── remove-duplicates_test.go # Тесты удаления дубликатов
│   └── parentheses_test.go    # Тесты проверки скобок
└── main.go                    # Точка входа для демонстрации
```

## 🚀 Запуск тестов

```bash
# Все тесты
cd tests && go test -v

# Конкретный тест
go test -v palindrome_test.go
go test -v backspace_test.go
go test -v reverse-string_test.go
go test -v remove-stars_test.go
go test -v remove-duplicates_test.go
go test -v parentheses_test.go

# Или используйте локальный скрипт из корня проекта
./scripts/check.sh
```

## 📋 Реализованные задачи

### Строки (Strings)
- **Palindrome** - проверка является ли строка палиндромом (LeetCode #125)
- **Reverse String** - разворот строки на месте (LeetCode #344)  
- **Backspace Compare** - сравнение строк с учетом backspace (LeetCode #844)
- **Remove Stars** - удаление звезд из строки (LeetCode #2390)
- **Remove Duplicates** - удаление дубликатов (LeetCode #1047)
- **Valid Parentheses** - проверка валидности скобок (LeetCode #20)

## 🎯 Особенности реализации

- Все функции оптимизированы по производительности
- Покрыты comprehensive тестами
- Следуют best practices Go
- Содержат подробные комментарии
- Автоматическое тестирование с CI/CD
- Бенчмарки для анализа производительности

## 📝 Формат задач

Каждая задача содержит:
- Описание проблемы и примеры
- Ограничения (constraints)
- Оптимизированное решение
- Comprehensive тесты

## 🔄 CI/CD

Проект использует GitHub Actions для автоматической проверки:
- **Тесты**: Автоматический запуск всех тестов при каждом пуше
- **Покрытие**: Анализ покрытия кода тестами
- **Бенчмарки**: Performance тестирование
- **Форматирование**: Проверка `go fmt`
- **Статический анализ**: Запуск `go vet`

Статус CI/CD: [![Go Tests](https://github.com/utkonoser/interviews/workflows/Go%20Tests/badge.svg)](https://github.com/utkonoser/interviews/actions/workflows/go-tests.yml)

## 🚀 Локальная разработка

Перед пушем запустите все проверки:
```bash
# Из корня проекта
./scripts/check.sh

# Или вручную
go test -v ./tests/...
go test -bench=. -benchmem ./tests/...
go fmt ./strings/...
go vet ./strings/...
```

---
**Удачи на собеседовании! 🚀**
