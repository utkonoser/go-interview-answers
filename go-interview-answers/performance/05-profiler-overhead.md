# 24. Overhead от стандартного профайлера?

## Ответ

**Стандартный профайлер Go имеет минимальный overhead благодаря семплированию.**

### Типы профилирования и их overhead:

#### 1. CPU Profile - минимальный overhead

```go
// CPU профилирование использует семплирование
// Overhead: ~5-10% в худшем случае
func cpuIntensive() {
    for i := 0; i < 1000000; i++ {
        math.Sqrt(float64(i))
    }
}
```

#### 2. Memory Profile - умеренный overhead

```go
// Memory профилирование более затратно
// Overhead: ~10-20% при активном использовании
func memoryIntensive() {
    for i := 0; i < 10000; i++ {
        data := make([]byte, 1024)
        _ = data
    }
}
```

#### 3. Goroutine Profile - минимальный overhead

```go
// Goroutine профилирование очень легкое
// Overhead: ~1-2%
func goroutineIntensive() {
    for i := 0; i < 1000; i++ {
        go func() {
            time.Sleep(time.Second)
        }()
    }
}
```

### Почему overhead минимален:

#### 1. Семплирующий профайлер

```go
// Go использует семплирование, а не инструментацию
// Собирает данные каждые 10ms (по умолчанию)
// Это означает очень низкий overhead

// В отличие от инструментированных профайлеров:
// - Не добавляет код в каждую функцию
// - Не замедляет выполнение
// - Собирает статистику выборочно
```

#### 2. Настройка частоты семплирования

```go
import "runtime/pprof"

func customProfile() {
    // Можно настроить частоту семплирования
    // По умолчанию: 1000Hz (1000 раз в секунду)
    
    // Более частое семплирование = больше overhead
    // Менее частое семплирование = меньше точность
}
```

### Измерение overhead:

```go
func measureOverhead() {
    // Без профайлера
    start := time.Now()
    cpuIntensive()
    withoutProfiler := time.Since(start)
    
    // С профайлером
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    
    start = time.Now()
    cpuIntensive()
    withProfiler := time.Since(start)
    
    pprof.StopCPUProfile()
    f.Close()
    
    overhead := float64(withProfiler-withProfiler) / float64(withoutProfiler) * 100
    fmt.Printf("Overhead: %.2f%%\n", overhead)
}
```

### Практические рекомендации:

#### 1. В разработке - всегда включен

```go
func main() {
    // В разработке overhead не критичен
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    
    // Ваше приложение
}
```

#### 2. В продакшене - по необходимости

```go
func main() {
    if os.Getenv("ENABLE_PROFILER") == "true" {
        go func() {
            http.ListenAndServe(":6060", nil)
        }()
    }
    
    // Ваше приложение
}
```

#### 3. Селективное профилирование

```go
func selectiveProfiling() {
    // Профилируем только определенные части
    if os.Getenv("PROFILE_SECTION") == "critical" {
        pprof.StartCPUProfile(os.Stdout)
        defer pprof.StopCPUProfile()
    }
    
    // Критический код
    criticalFunction()
}
```

### Сравнение с другими профайлерами:

```go
// Стандартный Go профайлер
// Overhead: 5-20%
// Точность: Хорошая
// Простота: Высокая

// Инструментированные профайлеры
// Overhead: 50-200%
// Точность: Отличная
// Простота: Низкая

// Аппаратные профайлеры
// Overhead: 1-5%
// Точность: Отличная
// Простота: Средняя
```

### Лучшие практики:

1. **Используйте в разработке** - overhead приемлем
2. **В продакшене включайте по необходимости**
3. **Мониторьте производительность** при включенном профайлере
4. **Используйте селективное профилирование** для критических участков
5. **Настройте частоту семплирования** под ваши нужды 