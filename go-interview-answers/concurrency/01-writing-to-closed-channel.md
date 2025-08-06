# 9. Что будет, если писать в закрытый канал?

## Ответ

**При записи в закрытый канал происходит паника (panic).** Это одна из немногих ситуаций в Go, когда канал может вызвать панику.

### Синтаксис записи в канал

```go
channel <- value
// Если канал закрыт, это вызовет panic
```

### Примеры

#### Базовый пример паники
```go
package main

import "fmt"

func main() {
    ch := make(chan int, 2)
    
    // Записываем данные
    ch <- 1
    ch <- 2
    
    // Закрываем канал
    close(ch)
    
    // Попытка записи в закрытый канал вызовет panic
    ch <- 3 // panic: send on closed channel
}
```

#### Безопасная запись с проверкой
```go
func safeWrite() {
    ch := make(chan int, 2)
    ch <- 1
    close(ch)
    
    // Безопасная запись с recover
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Паника перехвачена: %v\n", r)
        }
    }()
    
    ch <- 2 // Это вызовет panic
}
```

### Обработка паники

```go
func handlePanic() {
    ch := make(chan int)
    close(ch)
    
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Перехвачена паника: %v\n", r)
        }
    }()
    
    // Это вызовет panic
    ch <- 1
    fmt.Println("Эта строка не выполнится")
}
```

### Проверка состояния канала

```go
func checkChannelState() {
    ch := make(chan int, 1)
    
    // Проверяем, можно ли записать в канал
    select {
    case ch <- 1:
        fmt.Println("Успешно записали в канал")
    default:
        fmt.Println("Канал не готов к записи")
    }
    
    close(ch)
    
    // После закрытия канала
    select {
    case ch <- 2:
        fmt.Println("Это не выполнится")
    default:
        fmt.Println("Канал закрыт, запись невозможна")
    }
}
```

### Практические примеры

#### Graceful shutdown с проверкой
```go
func gracefulShutdownWithCheck() {
    ch := make(chan int, 5)
    done := make(chan bool)
    
    // Отправитель
    go func() {
        defer close(done)
        
        for i := 0; i < 10; i++ {
            select {
            case ch <- i:
                fmt.Printf("Отправлено: %d\n", i)
            case <-time.After(100 * time.Millisecond):
                fmt.Println("Таймаут отправки")
                return
            }
        }
    }()
    
    // Получатель
    go func() {
        for {
            select {
            case value, ok := <-ch:
                if !ok {
                    fmt.Println("Канал закрыт, завершаем работу")
                    return
                }
                fmt.Printf("Получено: %d\n", value)
            case <-time.After(200 * time.Millisecond):
                fmt.Println("Таймаут получения")
                close(ch)
                return
            }
        }
    }()
    
    <-done
}
```

#### Множественные отправители
```go
func multipleSenders() {
    ch := make(chan int, 10)
    done := make(chan bool)
    
    // Несколько отправителей
    for i := 0; i < 3; i++ {
        go func(id int) {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Printf("Отправитель %d: паника перехвачена: %v\n", id, r)
                }
            }()
            
            for j := 0; j < 5; j++ {
                select {
                case ch <- j:
                    fmt.Printf("Отправитель %d отправил: %d\n", id, j)
                case <-time.After(50 * time.Millisecond):
                    fmt.Printf("Отправитель %d: таймаут\n", id)
                    return
                }
            }
        }(i)
    }
    
    // Получатель
    go func() {
        count := 0
        for value := range ch {
            fmt.Printf("Получено: %d\n", value)
            count++
            if count >= 15 {
                close(ch)
                break
            }
        }
        close(done)
    }()
    
    <-done
}
```

### Сравнение с чтением

```go
func compareReadWrite() {
    ch := make(chan int, 1)
    ch <- 1
    close(ch)
    
    // Чтение из закрытого канала - безопасно
    value, ok := <-ch
    fmt.Printf("Чтение: %d, ok: %t\n", value, ok) // 1, true
    
    value, ok = <-ch
    fmt.Printf("Чтение после закрытия: %d, ok: %t\n", value, ok) // 0, false
    
    // Запись в закрытый канал - паника
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Паника при записи: %v\n", r)
        }
    }()
    
    ch <- 2 // panic: send on closed channel
}
```

### Лучшие практики

#### 1. Закрывайте канал только отправителем
```go
func bestPracticeSenderOnly() {
    ch := make(chan int)
    
    // Отправитель
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch) // Только отправитель закрывает канал
    }()
    
    // Получатель
    for value := range ch {
        fmt.Printf("Получено: %d\n", value)
    }
}
```

#### 2. Используйте select для безопасной записи
```go
func bestPracticeSelect() {
    ch := make(chan int, 1)
    close(ch)
    
    // Безопасная запись
    select {
    case ch <- 1:
        fmt.Println("Запись успешна")
    default:
        fmt.Println("Канал не готов к записи")
    }
}
```

#### 3. Используйте defer для обработки паники
```go
func bestPracticePanicHandling() {
    ch := make(chan int)
    close(ch)
    
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Обработка паники: %v\n", r)
        }
    }()
    
    // Потенциально опасная операция
    ch <- 1
}
```

#### 4. Проверяйте состояние канала
```go
func bestPracticeStateCheck() {
    ch := make(chan int, 1)
    
    // Проверяем готовность канала
    if len(ch) < cap(ch) {
        ch <- 1
        fmt.Println("Запись выполнена")
    } else {
        fmt.Println("Канал полон")
    }
    
    close(ch)
    
    // После закрытия проверяем через select
    select {
    case ch <- 2:
        fmt.Println("Это не выполнится")
    default:
        fmt.Println("Канал закрыт")
    }
}
```

### Дополнительные материалы

- [Go Tour: Channels](https://tour.golang.org/concurrency/2)
- [Effective Go: Channels](https://golang.org/doc/effective_go.html#channels)
- [Go by Example: Channels](https://gobyexample.com/channels)
- [Channel Axioms](https://dave.cheney.net/2014/03/19/channel-axioms) 