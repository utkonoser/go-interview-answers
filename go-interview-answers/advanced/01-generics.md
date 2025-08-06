# 26. Какие средства обобщенного программирования есть в Go?

## Ответ

**В Go есть несколько средств обобщенного программирования: дженерики (Go 1.18+), интерфейсы, кодогенерация, рефлексия.**

### 1. Дженерики (Go 1.18+)

```go
// Обобщенные функции
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Обобщенные типы
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}
```

### 2. Интерфейсы для полиморфизма

```go
// Интерфейсы как средство обобщения
type Processor interface {
    Process(data []byte) ([]byte, error)
}

type JSONProcessor struct{}
type XMLProcessor struct{}

func (j JSONProcessor) Process(data []byte) ([]byte, error) {
    // Обработка JSON
    return data, nil
}

func (x XMLProcessor) Process(data []byte) ([]byte, error) {
    // Обработка XML
    return data, nil
}

func processData(p Processor, data []byte) ([]byte, error) {
    return p.Process(data)
}
```

### 3. Кодогенерация

```go
//go:generate go run github.com/golang/mock/mockgen -source=interface.go -destination=mock.go

// Генерация кода для тестов
type UserService interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// Генерируется mock
// mockgen -source=service.go -destination=mock_service.go
```

### 4. Рефлексия

```go
import "reflect"

func MapSlice(slice interface{}, fn interface{}) interface{} {
    v := reflect.ValueOf(slice)
    fnv := reflect.ValueOf(fn)
    
    result := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
    
    for i := 0; i < v.Len(); i++ {
        args := []reflect.Value{v.Index(i)}
        result.Index(i).Set(fnv.Call(args)[0])
    }
    
    return result.Interface()
}

// Использование
numbers := []int{1, 2, 3, 4, 5}
doubled := MapSlice(numbers, func(x int) int {
    return x * 2
}).([]int)
```

### 5. Map-Reduce с дженериками

```go
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, item := range slice {
        result[i] = fn(item)
    }
    return result
}

func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, item := range slice {
        result = fn(result, item)
    }
    return result
}

// Использование
numbers := []int{1, 2, 3, 4, 5}
doubled := Map(numbers, func(x int) int { return x * 2 })
sum := Reduce(numbers, 0, func(acc, x int) int { return acc + x })
```

### 6. Обобщенные контейнеры

```go
// Обобщенная очередь
type Queue[T any] struct {
    items []T
}

func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}

// Обобщенный кэш
type Cache[K comparable, V any] struct {
    data map[K]V
    mu   sync.RWMutex
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, exists := c.data[key]
    return value, exists
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### 7. Ограничения типов

```go
// Числовые типы
func Sum[T constraints.Integer | constraints.Float](slice []T) T {
    var sum T
    for _, item := range slice {
        sum += item
    }
    return sum
}

// Сравнимые типы
func Contains[T comparable](slice []T, item T) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// Сложные ограничения
type Number interface {
    constraints.Integer | constraints.Float
}

func Abs[T Number](x T) T {
    if x < 0 {
        return -x
    }
    return x
}
```

### 8. Практические примеры

#### Обобщенный HTTP клиент

```go
type HTTPClient[T any] struct {
    client *http.Client
    baseURL string
}

func (h *HTTPClient[T]) Get(path string) (*T, error) {
    resp, err := h.client.Get(h.baseURL + path)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result T
    err = json.NewDecoder(resp.Body).Decode(&result)
    return &result, err
}

// Использование
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

client := &HTTPClient[User]{baseURL: "https://api.example.com"}
user, err := client.Get("/users/1")
```

#### Обобщенный валидатор

```go
type Validator[T any] struct {
    rules []func(T) error
}

func (v *Validator[T]) AddRule(rule func(T) error) {
    v.rules = append(v.rules, rule)
}

func (v *Validator[T]) Validate(item T) error {
    for _, rule := range v.rules {
        if err := rule(item); err != nil {
            return err
        }
    }
    return nil
}

// Использование
type User struct {
    Name  string
    Email string
    Age   int
}

validator := &Validator[User]{}
validator.AddRule(func(u User) error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    return nil
})
```

### Лучшие практики:

1. **Используйте дженерики** для обобщенных алгоритмов
2. **Интерфейсы** для полиморфизма
3. **Кодогенерация** для повторяющихся паттернов
4. **Рефлексию** только когда необходимо
5. **Ограничения типов** для безопасности 