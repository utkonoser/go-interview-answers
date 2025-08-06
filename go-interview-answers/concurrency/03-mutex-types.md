# 19. Какие типы мьютексов предоставляет stdlib?

## Ответ

**В Go есть два основных типа мьютексов: `sync.Mutex` и `sync.RWMutex`.**

### sync.Mutex - обычный мьютекс

```go
var mu sync.Mutex
var sharedData int

func writer() {
    mu.Lock()
    defer mu.Unlock()
    sharedData = 42
}

func reader() {
    mu.Lock()
    defer mu.Unlock()
    fmt.Println(sharedData)
}
```

### sync.RWMutex - мьютекс для чтения/записи

```go
var rwmu sync.RWMutex
var sharedData int

func writer() {
    rwmu.Lock()         // Эксклюзивная блокировка
    defer rwmu.Unlock()
    sharedData = 42
}

func reader() {
    rwmu.RLock()        // Блокировка для чтения
    defer rwmu.RUnlock()
    fmt.Println(sharedData)
}
```

### Когда использовать что:

#### sync.Mutex - когда:
- Много писателей
- Простая логика блокировки
- Производительность не критична

#### sync.RWMutex - когда:
- Много читателей, мало писателей
- Нужна оптимизация производительности
- Чтение не изменяет данные

### Практический пример:

```go
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### Дополнительные примитивы:

```go
// sync.Once - однократное выполнение
var once sync.Once
func init() {
    once.Do(func() {
        // Выполнится только один раз
    })
}

// sync.WaitGroup - ожидание завершения горутин
var wg sync.WaitGroup
func worker() {
    defer wg.Done()
    // Работа
}
wg.Add(1)
go worker()
wg.Wait()
``` 