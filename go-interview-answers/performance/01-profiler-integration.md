# 23. Как встроить стандартный профайлер в свое приложение?

## Ответ

**Стандартный профайлер Go встраивается через импорт пакета `net/http/pprof` и настройку HTTP-сервера.**

### Базовое встраивание

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof" // Импортируем pprof
)

func main() {
    // Запускаем HTTP-сервер на порту 6060
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Ваше приложение
    runYourApplication()
}

func runYourApplication() {
    // Симулируем нагрузку
    for i := 0; i < 1000000; i++ {
        // Какая-то работа
    }
}
```

### Доступные эндпоинты

После запуска профайлера доступны следующие URL:

- **CPU профиль**: `http://localhost:6060/debug/pprof/profile`
- **Heap профиль**: `http://localhost:6060/debug/pprof/heap`
- **Goroutine профиль**: `http://localhost:6060/debug/pprof/goroutine`
- **Block профиль**: `http://localhost:6060/debug/pprof/block`
- **Mutex профиль**: `http://localhost:6060/debug/pprof/mutex`
- **Thread профиль**: `http://localhost:6060/debug/pprof/threadcreate`
- **Общий профиль**: `http://localhost:6060/debug/pprof/`

### Программное включение профилирования

```go
package main

import (
    "log"
    "net/http"
    "runtime"
    "runtime/pprof"
    "time"
)

func main() {
    // Включаем профилирование CPU
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // Включаем профилирование памяти
    defer func() {
        memFile, err := os.Create("memory.prof")
        if err != nil {
            log.Fatal(err)
        }
        defer memFile.Close()
        
        runtime.GC() // Принудительная сборка мусора
        pprof.WriteHeapProfile(memFile)
    }()
    
    // Ваше приложение
    runYourApplication()
}

func runYourApplication() {
    // Симулируем нагрузку
    for i := 0; i < 1000000; i++ {
        // Какая-то работа
    }
}
```

### Настройка HTTP-сервера с профайлером

```go
package main

import (
    "log"
    "net/http"
    "net/http/pprof"
)

func main() {
    // Создаем мультиплексор
    mux := http.NewServeMux()
    
    // Добавляем эндпоинты профайлера
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    
    // Ваши эндпоинты
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    log.Println("Сервер запущен на :8080")
    log.Println("Профайлер доступен на :8080/debug/pprof/")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Профилирование конкретных участков кода

```go
package main

import (
    "log"
    "runtime/pprof"
    "time"
)

func main() {
    // Профилирование CPU для конкретной функции
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer cpuFile.Close()
    
    pprof.StartCPUProfile(cpuFile)
    defer pprof.StopCPUProfile()
    
    // Выполняем работу
    expensiveOperation()
}

func expensiveOperation() {
    // Симулируем дорогую операцию
    for i := 0; i < 1000000; i++ {
        // Какие-то вычисления
        _ = i * i
    }
}
```

### Профилирование памяти

```go
package main

import (
    "log"
    "runtime"
    "runtime/pprof"
)

func main() {
    // Включаем профилирование памяти
    memFile, err := os.Create("memory.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer memFile.Close()
    
    // Выполняем работу
    memoryIntensiveOperation()
    
    // Принудительная сборка мусора
    runtime.GC()
    
    // Записываем профиль памяти
    if err := pprof.WriteHeapProfile(memFile); err != nil {
        log.Fatal(err)
    }
}

func memoryIntensiveOperation() {
    // Симулируем интенсивное использование памяти
    var data []int
    for i := 0; i < 1000000; i++ {
        data = append(data, i)
    }
}
```

### Профилирование горутин

```go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "time"
)

func main() {
    // Запускаем HTTP-сервер для профайлера
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Создаем много горутин
    for i := 0; i < 1000; i++ {
        go func(id int) {
            for {
                time.Sleep(time.Second)
                // Какая-то работа
            }
        }(i)
    }
    
    // Ждем
    select {}
}
```

### Анализ профилей

#### Использование go tool pprof
```bash
# Анализ CPU профиля
go tool pprof cpu.prof

# Анализ памяти
go tool pprof memory.prof

# Анализ через HTTP
go tool pprof http://localhost:6060/debug/pprof/profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

#### Веб-интерфейс
```bash
# Запуск веб-интерфейса
go tool pprof -http=:8080 cpu.prof
```

### Программный анализ профилей

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "runtime/pprof"
)

func analyzeProfile() {
    // Получаем профиль CPU
    resp, err := http.Get("http://localhost:6060/debug/pprof/profile?seconds=30")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    // Анализируем профиль
    profile, err := pprof.NewProfileReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    // Получаем топ функций
    top := profile.Top()
    for i, sample := range top {
        fmt.Printf("%d. %s: %d\n", i+1, sample.Function, sample.Value)
    }
}
```

### Настройка для продакшена

```go
package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    // Проверяем переменную окружения
    if os.Getenv("ENABLE_PPROF") == "true" {
        go func() {
            log.Println("Профайлер включен на :6060")
            log.Fatal(http.ListenAndServe("localhost:6060", nil))
        }()
    }
    
    // Ваше приложение
    runApplication()
}

func runApplication() {
    // Основная логика приложения
}
```

### Мониторинг в реальном времени

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "runtime"
    "time"
)

type Stats struct {
    Goroutines int    `json:"goroutines"`
    Memory     uint64 `json:"memory"`
    Timestamp  string `json:"timestamp"`
}

func main() {
    // Эндпоинт для мониторинга
    http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        stats := Stats{
            Goroutines: runtime.NumGoroutine(),
            Memory:     m.Alloc,
            Timestamp:  time.Now().Format(time.RFC3339),
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
    })
    
    // Запускаем сервер
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()
    
    // Основное приложение
    runApplication()
}
```

### Лучшие практики

1. **Используйте отдельный порт** для профайлера
2. **Ограничивайте доступ** в продакшене
3. **Включайте профайлер** только при необходимости
4. **Анализируйте профили** регулярно
5. **Используйте веб-интерфейс** для удобства

### Дополнительные материалы

- [Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
- [Go pprof](https://pkg.go.dev/net/http/pprof)
- [runtime/pprof](https://pkg.go.dev/runtime/pprof) 