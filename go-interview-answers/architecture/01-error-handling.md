# 14. Предположим, ваша функция должна возвращать детализированные Recoverable и Fatal ошибки. Как это реализовано в пакете net? Как это надо делать в современном Go?

## Ответ

**В Go используется подход с возвращением ошибок как значений.** Начиная с Go 1.13, появились стандартные средства для работы с ошибками.

### Традиционный подход (до Go 1.13)

```go
package main

import (
    "errors"
    "fmt"
    "net"
)

// Определяем типы ошибок
var (
    ErrConnectionFailed = errors.New("connection failed")
    ErrTimeout         = errors.New("timeout")
    ErrInvalidData     = errors.New("invalid data")
)

// Функция, возвращающая ошибки
func connectToServer(host string, port int) error {
    if host == "" {
        return ErrInvalidData
    }
    
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return fmt.Errorf("failed to connect to %s: %w", address, err)
    }
    defer conn.Close()
    
    return nil
}
```

### Современный подход (Go 1.13+)

#### Использование errors.Is и errors.As
```go
package main

import (
    "errors"
    "fmt"
    "net"
    "time"
)

// Определяем типы ошибок
var (
    ErrConnectionFailed = errors.New("connection failed")
    ErrTimeout         = errors.New("timeout")
    ErrInvalidData     = errors.New("invalid data")
)

// Кастомная ошибка с дополнительной информацией
type ConnectionError struct {
    Host string
    Port int
    Err  error
}

func (e *ConnectionError) Error() string {
    return fmt.Sprintf("connection to %s:%d failed: %v", e.Host, e.Port, e.Err)
}

func (e *ConnectionError) Unwrap() error {
    return e.Err
}

// Функция с детализированными ошибками
func connectToServer(host string, port int) error {
    if host == "" {
        return fmt.Errorf("invalid host: %w", ErrInvalidData)
    }
    
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        return &ConnectionError{
            Host: host,
            Port: port,
            Err:  err,
        }
    }
    defer conn.Close()
    
    return nil
}
```

### Обработка ошибок

#### Проверка типа ошибки
```go
func handleErrors() {
    err := connectToServer("", 8080)
    
    // Проверяем конкретную ошибку
    if errors.Is(err, ErrInvalidData) {
        fmt.Println("Обрабатываем ошибку неверных данных")
        return
    }
    
    // Проверяем тип ошибки
    var connErr *ConnectionError
    if errors.As(err, &connErr) {
        fmt.Printf("Ошибка соединения с %s:%d: %v\n", 
            connErr.Host, connErr.Port, connErr.Err)
        return
    }
    
    // Проверяем системные ошибки
    if errors.Is(err, net.ErrClosed) {
        fmt.Println("Соединение закрыто")
        return
    }
}
```

### Иерархия ошибок

```go
// Базовые типы ошибок
type ErrorType int

const (
    ErrorTypeRecoverable ErrorType = iota
    ErrorTypeFatal
    ErrorTypeTemporary
)

// Структура для детализированных ошибок
type DetailedError struct {
    Type    ErrorType
    Message string
    Cause   error
    Context map[string]interface{}
}

func (e *DetailedError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *DetailedError) Unwrap() error {
    return e.Cause
}

// Конструкторы для разных типов ошибок
func NewRecoverableError(message string, cause error) *DetailedError {
    return &DetailedError{
        Type:    ErrorTypeRecoverable,
        Message: message,
        Cause:   cause,
    }
}

func NewFatalError(message string, cause error) *DetailedError {
    return &DetailedError{
        Type:    ErrorTypeFatal,
        Message: message,
        Cause:   cause,
    }
}

func NewTemporaryError(message string, cause error) *DetailedError {
    return &DetailedError{
        Type:    ErrorTypeTemporary,
        Message: message,
        Cause:   cause,
    }
}
```

### Практический пример

```go
func processData(data []byte) error {
    if len(data) == 0 {
        return NewRecoverableError("empty data provided", nil)
    }
    
    if len(data) > 1024 {
        return NewFatalError("data too large", nil)
    }
    
    // Симулируем обработку
    if data[0] == 0 {
        return NewTemporaryError("invalid data format", 
            errors.New("first byte cannot be zero"))
    }
    
    return nil
}

func handleProcessingError(err error) {
    var detailedErr *DetailedError
    if errors.As(err, &detailedErr) {
        switch detailedErr.Type {
        case ErrorTypeRecoverable:
            fmt.Printf("Восстанавливаемая ошибка: %s\n", detailedErr.Message)
            // Попробовать снова или использовать fallback
            
        case ErrorTypeFatal:
            fmt.Printf("Критическая ошибка: %s\n", detailedErr.Message)
            // Завершить работу
            
        case ErrorTypeTemporary:
            fmt.Printf("Временная ошибка: %s\n", detailedErr.Message)
            // Подождать и попробовать снова
        }
    } else {
        fmt.Printf("Неизвестная ошибка: %v\n", err)
    }
}
```

### Использование в пакете net

```go
func demonstrateNetErrors() {
    // Попытка подключения к несуществующему порту
    _, err := net.Dial("tcp", "localhost:99999")
    if err != nil {
        // Проверяем конкретные ошибки сети
        if errors.Is(err, net.ErrClosed) {
            fmt.Println("Соединение закрыто")
        } else if errors.Is(err, net.ErrInvalidAddr) {
            fmt.Println("Неверный адрес")
        } else {
            fmt.Printf("Ошибка сети: %v\n", err)
        }
    }
}
```

### Логирование ошибок

```go
import (
    "log"
    "runtime"
)

func logError(err error) {
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        log.Printf("Error at %s:%d: %v", file, line, err)
    }
}

func processWithLogging(data []byte) error {
    err := processData(data)
    if err != nil {
        logError(err)
    }
    return err
}
```

### Middleware для обработки ошибок

```go
type ErrorHandler func(error) error

func withErrorHandling(handler ErrorHandler) func([]byte) error {
    return func(data []byte) error {
        err := processData(data)
        if err != nil {
            return handler(err)
        }
        return nil
    }
}

func defaultErrorHandler(err error) error {
    var detailedErr *DetailedError
    if errors.As(err, &detailedErr) {
        // Добавляем контекст
        detailedErr.Context = map[string]interface{}{
            "timestamp": time.Now(),
            "retry_count": 0,
        }
    }
    return err
}
```

### Тестирование ошибок

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name    string
        data    []byte
        wantErr bool
        errType ErrorType
    }{
        {"empty data", []byte{}, true, ErrorTypeRecoverable},
        {"large data", make([]byte, 1025), true, ErrorTypeFatal},
        {"valid data", []byte{1, 2, 3}, false, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := processData(tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("processData() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if tt.wantErr {
                var detailedErr *DetailedError
                if errors.As(err, &detailedErr) {
                    if detailedErr.Type != tt.errType {
                        t.Errorf("expected error type %v, got %v", tt.errType, detailedErr.Type)
                    }
                }
            }
        })
    }
}
```

### Лучшие практики

1. **Используйте errors.Is** для проверки конкретных ошибок
2. **Используйте errors.As** для проверки типа ошибки
3. **Добавляйте контекст** с помощью fmt.Errorf и %w
4. **Создавайте иерархию ошибок** для сложных систем
5. **Логируйте ошибки** с достаточным контекстом
6. **Тестируйте обработку ошибок**

### Дополнительные материалы

- [Error handling and Go](https://blog.golang.org/error-handling-and-go)
- [Working with Errors in Go 1.13](https://blog.golang.org/go1.13-errors)
- [Effective Go: Errors](https://golang.org/doc/effective_go.html#errors) 