# Go Interview Tasks

Практические задачи по программированию на Go для подготовки к техническим собеседованиям.

## 📁 Структура

```
go-interview-tasks/
├── strings/                    # Задачи по работе со строками
│   ├── palindrome.go          # Проверка палиндрома
│   ├── reverse-string.go      # Разворот строки in-place
│   └── backspace.go           # Сравнение строк с backspace
├── tests/                     # Тесты для всех задач
│   ├── palindrome_test.go     # Тесты палиндрома
│   ├── reverse-string_test.go # Тесты разворота строки
│   └── backspace_test.go      # Тесты backspace
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
```

## 📋 Реализованные задачи

### Строки (Strings)
- **Palindrome** - проверка является ли строка палиндромом (LeetCode #125)
- **Reverse String** - разворот строки на месте (LeetCode #344)  
- **Backspace Compare** - сравнение строк с учетом backspace (LeetCode #844)

## 🎯 Особенности реализации

- Все функции оптимизированы по производительности
- Покрыты comprehensive тестами
- Следуют best practices Go
- Содержат подробные комментарии

## 📝 Формат задач

Каждая задача содержит:
- Описание проблемы и примеры
- Ограничения (constraints)
- Оптимизированное решение
- Comprehensive тесты

---
**Удачи на собеседовании! 🎉**
