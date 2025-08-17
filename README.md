# Go Interview Preparation

[![Go Tests](https://github.com/utkonoser/interviews/workflows/Go%20Tests/badge.svg)](https://github.com/utkonoser/interviews/actions/workflows/go-tests.yml)
[![Go Lint](https://github.com/utkonoser/interviews/workflows/Go%20Lint/badge.svg)](https://github.com/utkonoser/interviews/actions/workflows/go-lint.yml)

Полный набор материалов для подготовки к собеседованиям по Go: теоретические вопросы с ответами и практические задачи.

## 📁 Структура проекта

```
interviews/
├── go-interview-answers/       # Теоретические вопросы и ответы
│   ├── basics/                # Основы языка
│   ├── advanced/              # Продвинутые темы
│   ├── architecture/          # Архитектурные вопросы
│   ├── concurrency/           # Конкурентность
│   ├── performance/           # Производительность
│   └── README.md              # Полный список вопросов
├── go-interview-tasks/         # Практические задачи
│   ├── strings/               # Задачи по строкам
│   ├── tests/                 # Тесты для задач
│   └── README.md              # Описание задач
├── go.mod                     # Go module
└── LICENSE                    # Лицензия
```

## 🎯 Что включено

### 📚 Теоретические вопросы (28 вопросов)
- **Основы**: типы данных, слайсы, мапы, каналы, интерфейсы
- **Архитектура**: обработка ошибок, размещение интерфейсов, логирование
- **Конкурентность**: горутины, мьютексы, lock-free структуры  
- **Производительность**: профилирование, метрики, оптимизация
- **Продвинутые темы**: дженерики, преимущества/недостатки Go

### 💻 Практические задачи
- **Palindrome** - проверка палиндрома с оптимизацией
- **Reverse String** - разворот строки in-place
- **Backspace Compare** - сравнение строк с backspace

## 🚀 Быстрый старт

```bash
# Клонируйте репозиторий
git clone <repository-url>
cd interviews

# Изучите теоретические вопросы
open go-interview-answers/README.md

# Запустите тесты практических задач
cd go-interview-tasks/tests
go test -v
```

## 🔄 CI/CD

Проект использует GitHub Actions для автоматической проверки:

- **Тесты**: Автоматический запуск всех тестов при каждом пуше
- **Линтинг**: Проверка качества кода с golangci-lint
- **Покрытие**: Анализ покрытия кода тестами
- **Бенчмарки**: Performance тестирование
- **Множественные версии Go**: Тестирование на Go 1.21-1.24

Статус CI/CD отображается в README и на странице Actions.

## 📖 Как использовать

1. **Теория**: Изучите вопросы в `go-interview-answers/` - начните с `basics/`
2. **Практика**: Решите задачи в `go-interview-tasks/` и запустите тесты
3. **Подготовка**: Повторите сложные темы перед собеседованием

## 🎉 Особенности

- ✅ Полные ответы с примерами кода
- ✅ Оптимизированные решения задач  
- ✅ Comprehensive тесты
- ✅ Best practices Go
- ✅ Готово к использованию
- ✅ CI/CD с GitHub Actions
- ✅ Автоматическая проверка качества кода

---
**Удачи на собеседовании! 🚀**
