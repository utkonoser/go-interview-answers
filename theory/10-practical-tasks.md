# Практические задачи и алгоритмы

## 72. Реализовать LRU кэш

LRU (Least Recently Used) кэш — структура данных с ограниченной емкостью, удаляющая наименее используемые элементы. Реализация: комбинация hash map (O(1) доступ) и двусвязного списка (O(1) обновление порядка). При доступе элемент перемещается в голову списка. При переполнении удаляется хвост списка. Операции Get/Put за O(1).

```go
type LRUCache struct {
    capacity int
    cache    map[int]*Node // ключ → узел списка, O(1) поиск
    head     *Node         // самый «свежий» элемент
    tail     *Node         // самый «старый», его выкидываем при переполнении
}

type Node struct {
    key, value int
    prev, next *Node // двусвязный список для порядка использования
}

func (lru *LRUCache) Get(key int) int {
    if node, ok := lru.cache[key]; ok {
        lru.moveToHead(node) // прочитали — элемент стал недавно использованным
        return node.value
    }
    return -1 // ключ не найден
}

func (lru *LRUCache) Put(key, value int) {
    if node, ok := lru.cache[key]; ok {
        // ключ уже есть — обновляем значение и поднимаем в голову
        node.value = value
        lru.moveToHead(node)
    } else {
        if len(lru.cache) >= lru.capacity {
            lru.removeTail() // места нет — удаляем LRU-элемент с хвоста
        }
        newNode := &Node{key: key, value: value}
        lru.cache[key] = newNode
        lru.addToHead(newNode) // новый элемент — самый свежий
    }
}
// moveToHead, removeTail, addToHead — вспомогательные методы работы со списком
```

## 73. Реализовать thread-safe счетчик

Thread-safe счетчик с использованием `sync.Mutex` или `atomic` операций. Для простых операций `atomic` эффективнее, для сложной логики — мьютекс.

```go
// Вариант 1: atomic — быстрее для простого инкремента, без блокировок
type Counter struct {
    value int64
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.value, 1) // атомарно +1 на уровне CPU
}

func (c *Counter) Value() int64 {
    return atomic.LoadInt64(&c.value) // атомарное чтение
}

// Вариант 2: mutex — если логика сложнее одной операции
type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock() // разблокируем при выходе из функции
    c.value++
}

func (c *Counter) Value() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

## 74. Реализовать пул горутин (worker pool)

Worker pool — фиксированное количество worker-горутин обрабатывает задачи из канала. Контроль конкурентности, переиспользование горутин.

```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup

    // поднимаем фиксированное число воркеров — не по одной горутине на задачу
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs { // читаем, пока jobs не закроют
                results <- process(job)
            }
        }()
    }

    wg.Wait()      // ждём, пока все воркеры обработают очередь
    close(results) // после этого в results больше ничего не придёт
}
```

## 75. Реализовать rate limiter

Rate limiter ограничивает количество запросов за период времени. Реализации: token bucket, sliding window, fixed window.

```go
type RateLimiter struct {
    tokens chan struct{} // «жетоны»: забрал — разрешил запрос
    ticker *time.Ticker  // периодически пополняет bucket
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate), // буфер = максимум жетонов
        ticker: time.NewTicker(per / time.Duration(rate)),
    }

    // фоновая горутина: равномерно кидает жетоны в bucket
    go func() {
        for range rl.ticker.C {
            select {
            case rl.tokens <- struct{}{}: // положили жетон
            default: // bucket полон — лишний жетон выбрасываем
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens: // есть жетон — запрос разрешён
        return true
    default: // жетонов нет — лимит исчерпан
        return false
    }
}
```

## 76. Реализовать простой HTTP сервер с middleware

HTTP сервер с цепочкой middleware для логирования, аутентификации, обработки ошибок.

```go
// Middleware оборачивает handler и может сделать что-то до/после него
type Middleware func(http.HandlerFunc) http.HandlerFunc

func chain(middlewares ...Middleware) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        // оборачиваем с конца: первый в списке — внешний слой
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}

func logging(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r) // вызываем следующий handler в цепочке
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}

func main() {
    handler := chain(logging)(myHandler) // logging → myHandler
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

## 77. Реализовать producer-consumer паттерн

Producer генерирует данные, consumer обрабатывает. Синхронизация через каналы.

```go
func producer(items chan<- int) {
    for i := 0; i < 10; i++ {
        items <- i // кладём задачи в канал
    }
    close(items) // сигнал consumer'у: данных больше не будет
}

func consumer(items <-chan int, done chan<- bool) {
    for item := range items { // выходим, когда producer закроет канал
        process(item)
    }
    done <- true // сообщаем main, что всё обработано
}

func main() {
    items := make(chan int)
    done := make(chan bool)

    go producer(items)
    go consumer(items, done)

    <-done // main ждёт завершения consumer
}
```

## 78. Реализовать merge нескольких каналов

Объединение нескольких каналов в один с использованием `select` и горутин.

```go
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    // на каждый входной канал — своя горутина-перекачка
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v // пересылаем значения в общий канал
            }
        }(ch) // ch копируем в аргумент — иначе замыкание на последний ch
    }

    // закрываем out только когда все входные каналы исчерпаны
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

## 79. Реализовать timeout для операции

Таймаут для операции через `context.WithTimeout` или `select` с `time.After`.

```go
func withTimeout(fn func() error, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel() // освобождаем таймер, даже если fn уже вернулась

    done := make(chan error, 1) // буфер 1 — горутина не зависнет после return
    go func() {
        done <- fn() // тяжёлую работу делаем в отдельной горутине
    }()

    select {
    case err := <-done:
        return err // успели до дедлайна
    case <-ctx.Done():
        return ctx.Err() // context.DeadlineExceeded
    }
}
```

## 80. Реализовать graceful shutdown сервера

Graceful shutdown — корректное завершение с ожиданием завершения текущих запросов.

```go
func main() {
    server := &http.Server{Addr: ":8080"}

    // сервер в фоне — main свободен ловить сигналы ОС
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // Ctrl+C или kill
    <-quit // блокируемся до сигнала остановки

    // даём in-flight запросам до 5 секунд на завершение
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
}
```
