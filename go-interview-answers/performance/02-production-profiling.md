# 21. Способы поиска проблем производительности на проде?

## Ответ

**Основные способы: профилирование, метрики, логирование, мониторинг системы.**

### 1. Встроенный профайлер Go

```go
import (
    "net/http"
    _ "net/http/pprof"
)

func main() {
    // Включаем профайлер
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // Ваше приложение
}
```

### 2. Метрики Prometheus

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
        },
        []string{"method", "endpoint"},
    )
    
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
        },
        []string{"method", "endpoint", "status"},
    )
)

func init() {
    prometheus.MustRegister(requestDuration, requestCount)
}
```

### 3. Runtime метрики

```go
import (
    "runtime"
    "time"
)

func collectMetrics() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        // Логируем метрики
        log.Printf("Heap: %d MB, Goroutines: %d", 
            m.Alloc/1024/1024, runtime.NumGoroutine())
    }
}
```

### 4. Структурированное логирование

```go
import "github.com/sirupsen/logrus"

func handleRequest(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Обработка запроса
    processRequest(r)
    
    // Логируем производительность
    logrus.WithFields(logrus.Fields{
        "duration": time.Since(start),
        "method":   r.Method,
        "path":     r.URL.Path,
    }).Info("Request processed")
}
```

### 5. Системный мониторинг

```bash
# Мониторинг CPU и памяти
top -p $(pgrep your_app)

# Мониторинг сети
netstat -i

# Мониторинг диска
iostat -x 1

# Мониторинг процессов
ps aux | grep your_app
```

### 6. Анализ профилей

```bash
# CPU профиль
curl http://localhost:6060/debug/pprof/profile > cpu.prof

# Memory профиль
curl http://localhost:6060/debug/pprof/heap > mem.prof

# Goroutine профиль
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof
```

### 7. Инструменты анализа:

- **go tool pprof** - анализ профилей
- **Grafana** - визуализация метрик
- **Jaeger** - трейсинг запросов
- **Prometheus** - сбор метрик

### Лучшие практики:

1. **Всегда включайте профайлер** в продакшене
2. **Собирайте метрики** постоянно
3. **Настройте алерты** на критические метрики
4. **Используйте трейсинг** для сложных запросов
5. **Мониторьте внешние зависимости** 