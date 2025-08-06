# 25. Почему встраивание — не наследование?

## Ответ

**Встраивание в Go — это композиция, а не наследование.** Это принципиально разные концепции.

### Ключевые различия:

#### 1. Композиция vs Наследование

```go
// Встраивание (композиция)
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "Some sound"
}

type Dog struct {
    Animal // Встраивание
}

func (d Dog) Speak() string {
    return "Woof!" // Переопределение метода
}

// Использование
dog := Dog{Animal{Name: "Rex"}}
fmt.Println(dog.Name)    // Доступ к полю
fmt.Println(dog.Speak()) // Вызов переопределенного метода
```

#### 2. Отсутствие иерархии типов

```go
// В Go нет иерархии типов
type Animal struct{}
type Dog struct{ Animal }

func processAnimal(a Animal) {
    // Эта функция принимает только Animal
}

func main() {
    dog := Dog{}
    // processAnimal(dog) // ❌ Ошибка компиляции
    // Dog не является Animal в смысле наследования
}
```

#### 3. Интерфейсы для полиморфизма

```go
// В Go полиморфизм через интерфейсы
type Speaker interface {
    Speak() string
}

func makeSpeak(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    dog := Dog{Animal{Name: "Rex"}}
    makeSpeak(dog) // ✅ Работает через интерфейс
}
```

### Практические примеры:

#### 1. Встраивание структур

```go
type User struct {
    ID   int
    Name string
}

type Admin struct {
    User        // Встраивание
    Permissions []string
}

func (a Admin) IsAdmin() bool {
    return true
}

// Использование
admin := Admin{
    User: User{ID: 1, Name: "Admin"},
    Permissions: []string{"read", "write", "delete"},
}

fmt.Println(admin.Name)        // Доступ к полю User
fmt.Println(admin.IsAdmin())   // Метод Admin
```

#### 2. Встраивание интерфейсов

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader // Встраивание интерфейса
    Writer
}

// Реализация
type File struct {
    // Реализация методов
}

func (f File) Read(p []byte) (n int, err error) {
    return 0, nil
}

func (f File) Write(p []byte) (n int, err error) {
    return 0, nil
}

// File автоматически реализует ReadWriter
```

#### 3. Переопределение методов

```go
type Base struct {
    Name string
}

func (b Base) GetName() string {
    return b.Name
}

type Derived struct {
    Base
}

func (d Derived) GetName() string {
    return "Derived: " + d.Base.GetName() // Вызов базового метода
}
```

### Преимущества встраивания:

#### 1. Композиция лучше наследования

```go
// Гибкость композиции
type Logger struct {
    Level string
}

func (l Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.Level, msg)
}

type Database struct {
    Logger // Встраивание логгера
    URL    string
}

type HTTPClient struct {
    Logger // Тот же логгер
    Timeout time.Duration
}

// Один логгер используется в разных контекстах
```

#### 2. Избежание проблем наследования

```go
// Нет проблем с множественным наследованием
type A struct {
    Name string
}

type B struct {
    Name string
}

type C struct {
    A // Встраивание A
    B // Встраивание B
}

func main() {
    c := C{A: A{Name: "A"}, B: B{Name: "B"}}
    fmt.Println(c.A.Name) // Явный доступ к полю A
    fmt.Println(c.B.Name) // Явный доступ к полю B
}
```

### Сравнение с наследованием:

```go
// Наследование (в других языках)
class Animal {
    String name;
    void speak() { /* ... */ }
}

class Dog extends Animal {
    void speak() { /* ... */ } // Переопределение
}

// Встраивание (Go)
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "Some sound"
}

type Dog struct {
    Animal // Композиция
}

func (d Dog) Speak() string {
    return "Woof!" // Переопределение
}
```

### Лучшие практики:

#### 1. Используйте интерфейсы для полиморфизма

```go
type Processor interface {
    Process(data []byte) error
}

func processData(p Processor, data []byte) error {
    return p.Process(data)
}
```

#### 2. Композиция предпочтительнее наследования

```go
// ✅ Хорошо - композиция
type Service struct {
    Logger
    Database
    Cache
}

// ❌ Плохо - пытаться эмулировать наследование
```

#### 3. Явное переопределение методов

```go
type Base struct {
    Name string
}

func (b Base) GetName() string {
    return b.Name
}

type Derived struct {
    Base
}

func (d Derived) GetName() string {
    return "Derived: " + d.Base.GetName()
}
```

### Заключение:

- **Встраивание = композиция** - переиспользование кода
- **Интерфейсы = полиморфизм** - гибкость типов
- **Нет иерархии типов** - каждый тип независим
- **Явный доступ** - нет неоднозначностей 