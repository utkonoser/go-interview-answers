# 22. Стандартный набор метрик prometheus в Go-программе?

## Ответ

**Стандартные метрики включают: runtime метрики, HTTP метрики, бизнес-метрики.**

### 1. Runtime метрики (автоматически)

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// Эти метрики собираются автоматически
var (
    // Go runtime метрики
    goGoroutines = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "go_goroutines",
        Help: "Number of goroutines",
    })
    
    goThreads = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "go_threads",
        Help: "Number of OS threads",
    })
    
    goHeapAlloc = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "go_heap_alloc_bytes",
        Help: "Heap memory usage",
    })
)
```

### 2. HTTP метрики

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}
```

### 3. Бизнес-метрики

```go
var (
    // Метрики базы данных
    dbConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "db_connections_active",
        Help: "Active database connections",
    })
    
    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Database query duration",
        },
        []string{"query_type"},
    )
    
    // Метрики кэша
    cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_hits_total",
        Help: "Total cache hits",
    })
    
    cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cache_misses_total",
        Help: "Total cache misses",
    })
    
    // Метрики очередей
    queueSize = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "queue_size",
            Help: "Current queue size",
        },
        []string{"queue_name"},
    )
)
```

### 4. Middleware для HTTP метрик

```go
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Создаем response writer для захвата статуса
        rw := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(rw, r)
        
        // Записываем метрики
        duration := time.Since(start).Seconds()
        
        httpRequestsTotal.WithLabelValues(
            r.Method, 
            r.URL.Path, 
            strconv.Itoa(rw.statusCode),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            r.Method, 
            r.URL.Path,
        ).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 5. Экспорт метрик

```go
func main() {
    // Регистрируем метрики
    initMetrics()
    
    // Эндпоинт для метрик
    http.Handle("/metrics", promhttp.Handler())
    
    // Запускаем сервер
    http.ListenAndServe(":8080", nil)
}
```

### 6. Типичные метрики для мониторинга:

#### Системные:
- `go_goroutines` - количество горутин
- `go_heap_alloc_bytes` - использование памяти
- `go_gc_duration_seconds` - время сборки мусора

#### HTTP:
- `http_requests_total` - количество запросов
- `http_request_duration_seconds` - время ответа
- `http_requests_in_flight` - активные запросы

#### Бизнес:
- `user_registrations_total` - регистрации
- `payment_processed_total` - платежи
- `api_calls_total` - вызовы API

### 7. Конфигурация Prometheus:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'go-app'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Лучшие практики:

1. **Используйте осмысленные имена** метрик
2. **Добавляйте labels** для группировки
3. **Избегайте кардинальности** в labels
4. **Документируйте метрики** в help тексте
5. **Настройте алерты** на критические метрики 