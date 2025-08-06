# 20. Что такое lock-free структуры данных, и есть ли в Go такие?

## Ответ

**Lock-free структуры данных — это структуры, которые не используют блокировки для синхронизации.** В Go есть несколько встроенных lock-free примитивов.

### Встроенные lock-free примитивы:

#### 1. sync/atomic - атомарные операции

```go
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

func getValue() int64 {
    return atomic.LoadInt64(&counter)
}
```

#### 2. sync.Map - lock-free map

```go
var m sync.Map

func setValue(key, value string) {
    m.Store(key, value)
}

func getValue(key string) (string, bool) {
    value, ok := m.Load(key)
    if !ok {
        return "", false
    }
    return value.(string), true
}
```

### Lock-free структуры в Go:

#### 1. Каналы (частично lock-free)

```go
// Каналы используют lock-free алгоритмы внутри
ch := make(chan int, 10)

go func() {
    ch <- 42 // Lock-free отправка
}()

value := <-ch // Lock-free получение
```

#### 2. sync.Pool - lock-free пул

```go
var pool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func getBuffer() []byte {
    return pool.Get().([]byte)
}

func putBuffer(buf []byte) {
    pool.Put(buf)
}
```

### Создание собственных lock-free структур:

#### 1. Lock-free счетчик

```go
type LockFreeCounter struct {
    value int64
}

func (c *LockFreeCounter) Increment() {
    atomic.AddInt64(&c.value, 1)
}

func (c *LockFreeCounter) Get() int64 {
    return atomic.LoadInt64(&c.value)
}
```

#### 2. Lock-free стек (упрощенный)

```go
type Node struct {
    value interface{}
    next  *Node
}

type LockFreeStack struct {
    head *Node
}

func (s *LockFreeStack) Push(value interface{}) {
    newHead := &Node{value: value}
    for {
        oldHead := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.head)))
        newHead.next = (*Node)(oldHead)
        if atomic.CompareAndSwapPointer(
            (*unsafe.Pointer)(unsafe.Pointer(&s.head)),
            oldHead,
            unsafe.Pointer(newHead),
        ) {
            break
        }
    }
}
```

### Преимущества lock-free структур:

- **Производительность** - нет блокировок
- **Масштабируемость** - лучше на многоядерных системах
- **Отсутствие deadlock'ов** - нет блокировок

### Недостатки:

- **Сложность** - труднее отлаживать
- **Memory ordering** - сложности с памятью
- **ABA проблема** - в некоторых алгоритмах

### Когда использовать:

#### Lock-free - когда:
- Высокая конкуренция
- Критична производительность
- Простые операции

#### Обычные мьютексы - когда:
- Простота важнее производительности
- Сложная логика синхронизации
- Редкая конкуренция 