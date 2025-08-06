# 23. Как встроить стандартный профайлер в свое приложение?

## Ответ

**Используйте пакет `net/http/pprof` для встраивания профайлера.**

### 1. Базовое встраивание

```go
import (
    "net/http"
    _ "net/http/pprof" // Автоматически регистрирует эндпоинты
)

func main() {
    // Запускаем профайлер на отдельном порту
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // Ваше приложение
    http.ListenAndServe(":8080", nil)
}
```

### 2. Интеграция с основным сервером

```go
import (
    "net/http"
    "net/http/pprof"
)

func main() {
    mux := http.NewServeMux()
    
    // Основные эндпоинты
    mux.HandleFunc("/", handleHome)
    mux.HandleFunc("/api/users", handleUsers)
    
    // Профайлер эндпоинты
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    
    http.ListenAndServe(":8080", mux)
}
```

### 3. Программное профилирование

```go
import (
    "os"
    "runtime/pprof"
    "time"
)

func profileCPU() {
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Ваш код для профилирования
    time.Sleep(30 * time.Second)
}

func profileMemory() {
    f, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    pprof.WriteHeapProfile(f)
}
```

### 4. Middleware для профилирования

```go
func profilingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/debug/pprof/profile" {
            // Запускаем CPU профилирование
            pprof.Profile(w, r)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### 5. Конфигурация для продакшена

```go
func setupProfiler() {
    if os.Getenv("ENV") == "production" {
        // В продакшене ограничиваем доступ
        http.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
            // Проверяем авторизацию
            if !isAuthorized(r) {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            pprof.Index(w, r)
        })
    }
}
```

### 6. Анализ профилей

```bash
# CPU профиль
curl http://localhost:6060/debug/pprof/profile > cpu.prof
go tool pprof cpu.prof

# Memory профиль
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Goroutine профиль
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof

# Trace профиль
curl http://localhost:6060/debug/pprof/trace?seconds=30 > trace.out
go tool trace trace.out
```

### 7. Автоматическое профилирование

```go
func autoProfile() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            // CPU профиль каждые 30 секунд
            f, _ := os.Create(fmt.Sprintf("cpu_%d.prof", time.Now().Unix()))
            pprof.StartCPUProfile(f)
            time.Sleep(30 * time.Second)
            pprof.StopCPUProfile()
            f.Close()
            
            // Memory профиль
            mf, _ := os.Create(fmt.Sprintf("mem_%d.prof", time.Now().Unix()))
            pprof.WriteHeapProfile(mf)
            mf.Close()
        }
    }()
}
```

### Доступные эндпоинты:

- `/debug/pprof/` - главная страница
- `/debug/pprof/profile` - CPU профиль
- `/debug/pprof/heap` - Memory профиль
- `/debug/pprof/goroutine` - Goroutine профиль
- `/debug/pprof/block` - Block профиль
- `/debug/pprof/mutex` - Mutex профиль
- `/debug/pprof/trace` - Trace профиль

### Лучшие практики:

1. **Используйте отдельный порт** для профайлера
2. **Ограничьте доступ** в продакшене
3. **Настройте ротацию** профилей
4. **Мониторьте размер** профилей
5. **Используйте алерты** на аномалии 