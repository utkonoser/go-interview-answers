# 18. Можно ли использовать один и тот же буфер []byte в нескольких горутинах?

## Ответ

**Технически можно, но практически не нужно.** Использование общего буфера с мьютексом — антипаттерн.

### Проблемы общего буфера:

```go
// ❌ Плохо - общий буфер с мьютексом
var sharedBuffer []byte
var mu sync.Mutex

func worker1() {
    mu.Lock()
    defer mu.Unlock()
    copy(sharedBuffer, []byte("data1"))
}

func worker2() {
    mu.Lock()
    defer mu.Unlock()
    copy(sharedBuffer, []byte("data2"))
}
```

### Правильные подходы:

#### 1. Локальные буферы

```go
// ✅ Хорошо - каждый горутина свой буфер
func worker(id int) {
    buffer := make([]byte, 1024)
    // Работаем с локальным буфером
    copy(buffer, []byte(fmt.Sprintf("data%d", id)))
}
```

#### 2. Каналы для передачи данных

```go
// ✅ Хорошо - передача через каналы
func producer(ch chan []byte) {
    data := []byte("hello")
    ch <- data
}

func consumer(ch chan []byte) {
    data := <-ch
    // Обрабатываем полученные данные
}
```

#### 3. sync.Pool для переиспользования

```go
// ✅ Хорошо - пул буферов
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func worker() {
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer)
    
    // Используем буфер
    copy(buffer, []byte("data"))
}
```

### Когда действительно нужен общий буфер:

```go
// Редкий случай - кольцевой буфер
type RingBuffer struct {
    buffer []byte
    mu     sync.RWMutex
    read   int
    write  int
}

func (rb *RingBuffer) Write(data []byte) {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    // Запись в буфер
}

func (rb *RingBuffer) Read(size int) []byte {
    rb.mu.RLock()
    defer rb.mu.RUnlock()
    // Чтение из буфера
    return make([]byte, size) // Упрощенно
}
```

### Лучшие практики:

1. **Избегайте общего состояния** - используйте локальные переменные
2. **Передавайте данные через каналы** - идиоматично для Go
3. **Используйте sync.Pool** для переиспользования объектов
4. **Рассмотрите lock-free структуры** для высоконагруженных систем 