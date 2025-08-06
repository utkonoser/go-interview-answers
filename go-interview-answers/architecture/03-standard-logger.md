# 15. Главный недостаток стандартного логгера?

## Ответ

**Главный недостаток стандартного логгера Go — отсутствие уровней логирования и структурированного логирования.** Это делает его непригодным для продакшн-среды.

### Проблемы стандартного логгера

#### 1. Отсутствие уровней логирования

```go
package main

import (
    "log"
    "os"
)

func demonstrateStandardLogger() {
    // Стандартный логгер не имеет уровней
    log.Println("Это информационное сообщение")
    log.Println("Это ошибка")
    log.Println("Это предупреждение")
    
    // Нет возможности фильтровать по уровням
    // Нет возможности отключить определенные типы сообщений
}
```

#### 2. Отсутствие структурированного логирования

```go
package main

import (
    "log"
    "time"
)

func demonstrateUnstructuredLogging() {
    userID := 123
    action := "login"
    timestamp := time.Now()
    
    // Плохо: неструктурированный лог
    log.Printf("User %d performed %s at %v", userID, action, timestamp)
    
    // Нет возможности легко парсить логи
    // Нет возможности агрегировать по полям
    // Сложно анализировать в системах мониторинга
}
```

#### 3. Ограниченная настройка

```go
package main

import (
    "log"
    "os"
)

func demonstrateLimitedConfiguration() {
    // Минимальные возможности настройки
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    log.SetOutput(os.Stderr)
    
    // Нет возможности:
    // - Настроить формат времени
    // - Добавить контекстные поля
    // - Настроить ротацию логов
    // - Отправить логи в внешние системы
}
```

### Современные альтернативы

#### 1. logrus - структурированный логгер

```go
package main

import (
    "github.com/sirupsen/logrus"
    "time"
)

func demonstrateLogrus() {
    logger := logrus.New()
    
    // Настройка уровней
    logger.SetLevel(logrus.InfoLevel)
    
    // Структурированное логирование
    logger.WithFields(logrus.Fields{
        "user_id": 123,
        "action":  "login",
        "ip":      "192.168.1.1",
    }).Info("User logged in")
    
    // Логирование ошибок
    logger.WithFields(logrus.Fields{
        "error":   "database connection failed",
        "retries": 3,
    }).Error("Failed to connect to database")
    
    // Логирование с контекстом
    logger.WithField("request_id", "abc-123").Info("Request processed")
}
```

#### 2. zap - высокопроизводительный логгер

```go
package main

import (
    "go.uber.org/zap"
    "time"
)

func demonstrateZap() {
    // Создание логгера
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    // Структурированное логирование
    logger.Info("User logged in",
        zap.Int("user_id", 123),
        zap.String("action", "login"),
        zap.String("ip", "192.168.1.1"),
    )
    
    // Логирование ошибок
    logger.Error("Database connection failed",
        zap.String("error", "connection timeout"),
        zap.Int("retries", 3),
    )
    
    // Логирование с контекстом
    logger.Info("Request processed",
        zap.String("request_id", "abc-123"),
        zap.Duration("duration", time.Millisecond*150),
    )
}
```

#### 3. zerolog - JSON логгер

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "time"
)

func demonstrateZerolog() {
    // Настройка для продакшена
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    
    // Структурированное логирование
    log.Info().
        Int("user_id", 123).
        Str("action", "login").
        Str("ip", "192.168.1.1").
        Msg("User logged in")
    
    // Логирование ошибок
    log.Error().
        Str("error", "database connection failed").
        Int("retries", 3).
        Msg("Failed to connect to database")
    
    // Логирование с контекстом
    log.Info().
        Str("request_id", "abc-123").
        Dur("duration", time.Millisecond*150).
        Msg("Request processed")
}
```

### Практические примеры

#### 1. HTTP сервер с логированием

```go
package main

import (
    "net/http"
    "time"
    "github.com/sirupsen/logrus"
)

var logger = logrus.New()

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Логируем начало запроса
        logger.WithFields(logrus.Fields{
            "method":     r.Method,
            "path":       r.URL.Path,
            "remote_addr": r.RemoteAddr,
            "user_agent": r.UserAgent(),
        }).Info("Request started")
        
        // Обрабатываем запрос
        next.ServeHTTP(w, r)
        
        // Логируем завершение запроса
        logger.WithFields(logrus.Fields{
            "method":     r.Method,
            "path":       r.URL.Path,
            "duration":   time.Since(start),
            "status":     200, // В реальности получаем из response
        }).Info("Request completed")
    })
}

func main() {
    // Настройка логгера
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetLevel(logrus.InfoLevel)
    
    // Создаем сервер с middleware
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

#### 2. База данных с логированием

```go
package main

import (
    "database/sql"
    "github.com/sirupsen/logrus"
    _ "github.com/lib/pq"
)

type Database struct {
    db     *sql.DB
    logger *logrus.Logger
}

func NewDatabase(dsn string, logger *logrus.Logger) (*Database, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        logger.WithError(err).Error("Failed to open database connection")
        return nil, err
    }
    
    logger.Info("Database connection established")
    
    return &Database{
        db:     db,
        logger: logger,
    }, nil
}

func (d *Database) GetUser(id int) (*User, error) {
    d.logger.WithField("user_id", id).Debug("Fetching user from database")
    
    user := &User{}
    err := d.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).
        Scan(&user.ID, &user.Name, &user.Email)
    
    if err != nil {
        d.logger.WithFields(logrus.Fields{
            "user_id": id,
            "error":   err.Error(),
        }).Error("Failed to fetch user from database")
        return nil, err
    }
    
    d.logger.WithField("user_id", id).Debug("User fetched successfully")
    return user, nil
}
```

#### 3. Конфигурация логгера

```go
package main

import (
    "os"
    "github.com/sirupsen/logrus"
)

func setupLogger() *logrus.Logger {
    logger := logrus.New()
    
    // Настройка уровня логирования
    level := os.Getenv("LOG_LEVEL")
    switch level {
    case "debug":
        logger.SetLevel(logrus.DebugLevel)
    case "info":
        logger.SetLevel(logrus.InfoLevel)
    case "warn":
        logger.SetLevel(logrus.WarnLevel)
    case "error":
        logger.SetLevel(logrus.ErrorLevel)
    default:
        logger.SetLevel(logrus.InfoLevel)
    }
    
    // Настройка формата
    if os.Getenv("LOG_FORMAT") == "json" {
        logger.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }
    
    // Настройка вывода
    if os.Getenv("LOG_OUTPUT") == "file" {
        file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err == nil {
            logger.SetOutput(file)
        }
    }
    
    return logger
}
```

### Сравнение логгеров

```go
package main

import (
    "log"
    "github.com/sirupsen/logrus"
    "go.uber.org/zap"
    "github.com/rs/zerolog/log"
)

func compareLoggers() {
    // Стандартный логгер
    log.Println("Standard logger: basic message")
    
    // Logrus
    logrus.WithField("user_id", 123).Info("Logrus: structured message")
    
    // Zap
    logger, _ := zap.NewProduction()
    logger.Info("Zap: structured message", zap.Int("user_id", 123))
    
    // Zerolog
    log.Info().Int("user_id", 123).Msg("Zerolog: structured message")
}
```

### Лучшие практики

#### 1. Используйте структурированное логирование

```go
// Плохо
log.Printf("User %d performed %s", userID, action)

// Хорошо
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  action,
    "timestamp": time.Now(),
}).Info("User action performed")
```

#### 2. Настраивайте уровни логирования

```go
func setupLogLevels() {
    switch os.Getenv("ENV") {
    case "development":
        logger.SetLevel(logrus.DebugLevel)
    case "staging":
        logger.SetLevel(logrus.InfoLevel)
    case "production":
        logger.SetLevel(logrus.WarnLevel)
    }
}
```

#### 3. Добавляйте контекст

```go
func logWithContext(ctx context.Context, message string) {
    logger.WithFields(logrus.Fields{
        "request_id": ctx.Value("request_id"),
        "user_id":    ctx.Value("user_id"),
        "session_id": ctx.Value("session_id"),
    }).Info(message)
}
```

#### 4. Используйте правильные уровни

```go
func demonstrateLogLevels() {
    logger := logrus.New()
    
    // Debug - для отладки
    logger.Debug("Processing request step 1")
    
    // Info - для общей информации
    logger.Info("User logged in successfully")
    
    // Warn - для предупреждений
    logger.Warn("Database connection slow")
    
    // Error - для ошибок
    logger.Error("Failed to process payment")
    
    // Fatal - для критических ошибок
    logger.Fatal("Cannot start application")
}
```

### Дополнительные материалы

- [logrus](https://github.com/sirupsen/logrus)
- [zap](https://github.com/uber-go/zap)
- [zerolog](https://github.com/rs/zerolog)
- [Go Logging](https://golang.org/pkg/log/) 