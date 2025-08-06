# Go Interview Questions & Answers

Полный набор ответов на вопросы по Go для подготовки к собеседованиям.

## 📁 Структура проекта

```
go-interview-answers/
├── advanced/           # Продвинутые темы
│   ├── 01-generics.md                    # Дженерики
│   ├── 02-go-advantages-disadvantages.md # Преимущества и недостатки Go
│   ├── 03-embedding-vs-inheritance.md   # Встраивание vs наследование
│   └── 04-linters.md                    # Линтеры
├── architecture/       # Архитектура
│   ├── 01-error-handling.md             # Обработка ошибок
│   ├── 02-interface-placement.md        # Размещение интерфейсов
│   ├── 03-standard-logger.md            # Стандартный логгер
│   └── 04-orm-go.md                     # ORM в Go
├── basics/            # Основы
│   ├── 01-imperative-declarative.md     # Императивное vs декларативное
│   ├── 02-type-switch.md                # Type switch
│   ├── 03-interface-implementation.md   # Реализация интерфейсов
│   ├── 04-append-function.md            # Функция append
│   ├── 05-slice-zero-value.md           # Нулевое значение slice
│   ├── 06-map-structure.md              # Структура map
│   ├── 07-map-iteration-order.md        # Порядок итерации map
│   ├── 08-reading-from-closed-channel.md # Чтение из закрытого канала
│   ├── 10-sorting-structures.md         # Сортировка структур
│   ├── 11-serialization.md              # Сериализация
│   └── 12-linked-list-reversal.md       # Разворот связного списка
├── concurrency/       # Конкурентность
│   ├── 01-writing-to-closed-channel.md  # Запись в закрытый канал
│   ├── 02-buffer-sharing.md             # Общие буферы в горутинах
│   ├── 03-mutex-types.md                # Типы мьютексов
│   └── 04-lock-free-structures.md       # Lock-free структуры
└── performance/       # Производительность
    ├── 01-profiler-integration.md       # Интеграция профайлера
    ├── 02-production-profiling.md       # Профилирование на проде
    ├── 03-prometheus-metrics.md         # Метрики Prometheus
    ├── 04-profiler-integration.md       # Встраивание профайлера
    └── 05-profiler-overhead.md          # Overhead профайлера
```

## 📋 Полный список вопросов

### Основы (Basics)

1. **[Go — императивный или декларативный? А в чем разница?](go-interview-answers/basics/01-imperative-declarative.md)**
   - Императивное vs декларативное программирование
   - Примеры в Go

2. **[Что такое type switch?](go-interview-answers/basics/02-type-switch.md)**
   - Работа с типами в runtime
   - Type assertions и type switches

3. **[Как сообщить компилятору, что наш тип реализует интерфейс?](go-interview-answers/basics/03-interface-implementation.md)**
   - Duck typing в Go
   - Неявная реализация интерфейсов

4. **[Как работает append?](go-interview-answers/basics/04-append-function.md)**
   - Управление памятью в Go
   - Поведение append с capacity

5. **[Какое у slice zero value? Какие операции над ним возможны?](go-interview-answers/basics/05-slice-zero-value.md)**
   - Нулевое значение slice
   - Операции с nil slice

6. **[Как устроен тип map?](go-interview-answers/basics/06-map-structure.md)**
   - Внутренняя структура map
   - Hash таблицы в Go

7. **[Каков порядок перебора map?](go-interview-answers/basics/07-map-iteration-order.md)**
   - Случайный порядок итерации
   - Почему порядок не гарантирован

8. **[Что будет, если читать из закрытого канала?](go-interview-answers/basics/08-reading-from-closed-channel.md)**
   - Поведение закрытых каналов
   - Zero values для типов канала

9. **[Что будет, если писать в закрытый канал?](go-interview-answers/concurrency/01-writing-to-closed-channel.md)**
   - Panic при записи в закрытый канал
   - Правильное закрытие каналов

10. **[Как вы отсортируете массив структур по алфавиту по полю Name?](go-interview-answers/basics/10-sorting-structures.md)**
    - Сортировка структур
    - Использование sort.Slice

11. **[Что такое сериализация? Зачем она нужна?](go-interview-answers/basics/11-serialization.md)**
    - Сериализация и десериализация
    - JSON, XML, Protocol Buffers

12. **[Сколько времени в минутах займет у вас написание процедуры обращения односвязного списка?](go-interview-answers/basics/12-linked-list-reversal.md)**
    - Алгоритм разворота списка
    - Итеративный и рекурсивный подходы

### Архитектура (Architecture)

13. **[Где следует поместить описание интерфейса: в пакете с реализацией или в пакете, где этот интерфейс используется?](go-interview-answers/architecture/02-interface-placement.md)**
    - Принципы размещения интерфейсов
    - Dependency inversion

14. **[Предположим, ваша функция должна возвращать детализированные Recoverable и Fatal ошибки. Как это реализовано в пакете net?](go-interview-answers/architecture/01-error-handling.md)**
    - Обработка ошибок в Go
    - Wrapping и unwrapping ошибок

15. **[Главный недостаток стандартного логгера?](go-interview-answers/architecture/03-standard-logger.md)**
    - Проблемы стандартного логгера
    - Альтернативы (logrus, zap, zerolog)

16. **[Есть ли для Go хороший orm? Ответ обоснуйте.](go-interview-answers/architecture/04-orm-go.md)**
    - ORM vs SQL
    - GORM, SQLx, database/sql

### Продвинутые темы (Advanced)

17. **[Какой у вас любимый линтер?](go-interview-answers/advanced/04-linters.md)**
    - golangci-lint
    - Статический анализ кода

25. **[Почему встраивание — не наследование?](go-interview-answers/advanced/03-embedding-vs-inheritance.md)**
    - Композиция vs наследование
    - Встраивание структур в Go

26. **[Какие средства обобщенного программирования есть в Go?](go-interview-answers/advanced/01-generics.md)**
    - Дженерики (Go 1.18+)
    - Интерфейсы и рефлексия

27. **[Какие технологические преимущества языка Go вы можете назвать?](go-interview-answers/advanced/02-go-advantages-disadvantages.md)**
    - Простота и производительность
    - Конкурентность и сетевые возможности

28. **[Какие технологические недостатки языка Go вы можете назвать?](go-interview-answers/advanced/02-go-advantages-disadvantages.md)**
    - Ограничения системы типов
    - Обработка ошибок

### Конкурентность (Concurrency)

18. **[Можно ли использовать один и тот же буфер []byte в нескольких горутинах?](go-interview-answers/concurrency/02-buffer-sharing.md)**
    - Общие ресурсы в горутинах
    - Правильные подходы к синхронизации

19. **[Какие типы мьютексов предоставляет stdlib?](go-interview-answers/concurrency/03-mutex-types.md)**
    - sync.Mutex и sync.RWMutex
    - Примитивы синхронизации

20. **[Что такое lock-free структуры данных, и есть ли в Go такие?](go-interview-answers/concurrency/04-lock-free-structures.md)**
    - Атомарные операции
    - sync.Map и sync/atomic

### Производительность (Performance)

21. **[Способы поиска проблем производительности на проде?](go-interview-answers/performance/02-production-profiling.md)**
    - Профилирование в продакшене
    - Метрики и мониторинг

22. **[Стандартный набор метрик prometheus в Go-программе?](go-interview-answers/performance/03-prometheus-metrics.md)**
    - Runtime метрики
    - HTTP метрики и бизнес-метрики

23. **[Как встроить стандартный профайлер в свое приложение?](go-interview-answers/performance/04-profiler-integration.md)**
    - net/http/pprof
    - Программное профилирование

24. **[Overhead от стандартного профайлера?](go-interview-answers/performance/05-profiler-overhead.md)**
    - Семплирующий профайлер
    - Минимальный overhead

## 🎯 Темы

### Основы (Basics)
- Типы данных и их нулевые значения
- Слайсы и их поведение
- Мапы и их структура
- Каналы и их использование
- Интерфейсы и их реализация

### Архитектура (Architecture)
- Обработка ошибок в Go
- Размещение интерфейсов
- Логирование и стандартный логгер
- Работа с базами данных и ORM

### Конкурентность (Concurrency)
- Горутины и их использование
- Каналы и их типы
- Синхронизация с мьютексами
- Lock-free структуры данных

### Производительность (Performance)
- Профилирование приложений
- Метрики и мониторинг
- Оптимизация производительности
- Анализ профилей

### Продвинутые темы (Advanced)
- Дженерики (Go 1.18+)
- Преимущества и недостатки Go
- Встраивание vs наследование
- Линтеры и инструменты

## 📚 Как использовать

1. **Выберите вопрос** из списка выше
2. **Перейдите по ссылке** на соответствующий ответ
3. **Изучите ответ** с примерами кода
4. **Попрактикуйтесь** с предоставленными примерами
5. **Подготовьтесь** к вопросам на собеседовании

## 🚀 Быстрый старт

```bash
# Клонируйте репозиторий
git clone https://github.com/utkonoser/go-interview-answers.git

# Перейдите в папку
cd go-interview-answers

# Откройте нужный файл
open go-interview-answers/basics/01-imperative-declarative.md
```

## 📝 Формат ответов

Каждый ответ содержит:
- **Краткое объяснение** концепции
- **Практические примеры** кода
- **Лучшие практики**
- **Частые ошибки** и как их избежать
- **Дополнительные материалы**

## 🤝 Вклад в проект

Если вы нашли ошибку или хотите добавить что-то:

1. Создайте fork репозитория
2. Создайте ветку для ваших изменений
3. Внесите изменения
4. Создайте Pull Request

## 📄 Лицензия

MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🙏 Благодарности

Спасибо сообществу Go за отличную документацию и примеры!

---

**Удачи на собеседовании! 🎉** 