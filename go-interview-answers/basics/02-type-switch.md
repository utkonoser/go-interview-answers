# 2. Что такое type switch?

## Ответ

**Type switch** — это конструкция в Go для определения типа значения интерфейса во время выполнения.

### Синтаксис

```go
switch v := x.(type) {
case T1:
    // v имеет тип T1
case T2:
    // v имеет тип T2
default:
    // v имеет другой тип
}
```

### Базовый пример

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case string:
        fmt.Printf("Строка: %s\n", v)
    case int:
        fmt.Printf("Целое число: %d\n", v)
    case bool:
        fmt.Printf("Булево значение: %t\n", v)
    default:
        fmt.Printf("Неизвестный тип: %T\n", v)
    }
}

func main() {
    describe("hello")
    describe(42)
    describe(true)
    describe(3.14)
}
```

### Практический пример: Тип-сумма

```go
// Определяем тип-сумму для числовых значений
type Number interface {
    isNumber()
}

type Int64 int64
type Float64 float64
type Complex complex128

func (Int64) isNumber()   {}
func (Float64) isNumber() {}
func (Complex) isNumber() {}

// Реализация метода Add с использованием type switch
func (n Number) Add(other int64) Number {
    switch v := n.(type) {
    case Int64:
        return Int64(int64(v) + other)
    case Float64:
        return Float64(float64(v) + float64(other))
    case Complex:
        return Complex(complex(float64(real(complex128(v)))+float64(other), imag(complex128(v))))
    default:
        return Int64(other)
    }
}
```

### Использование с интерфейсами

```go
type Animal interface {
    Speak() string
}

type Dog struct{}
type Cat struct{}

func (d Dog) Speak() string { return "Гав!" }
func (c Cat) Speak() string { return "Мяу!" }

func processAnimal(a Animal) {
    switch v := a.(type) {
    case Dog:
        fmt.Printf("Собака говорит: %s\n", v.Speak())
    case Cat:
        fmt.Printf("Кошка говорит: %s\n", v.Speak())
    default:
        fmt.Printf("Неизвестное животное: %s\n", v.Speak())
    }
}
```

### Type switch vs Type assertion

#### Type assertion (простое приведение типа)
```go
value, ok := x.(string)
if ok {
    fmt.Println("Это строка:", value)
}
```

#### Type switch (множественное приведение)
```go
switch v := x.(type) {
case string:
    fmt.Println("Строка:", v)
case int:
    fmt.Println("Число:", v)
}
```

### Особенности

1. **Переменная `v`** в каждом case имеет соответствующий тип
2. **Fallthrough** не работает в type switch
3. **Default case** обрабатывает все остальные типы
4. **Можно использовать** для проверки nil

### Пример с nil

```go
func processValue(x interface{}) {
    switch v := x.(type) {
    case nil:
        fmt.Println("Значение равно nil")
    case string:
        fmt.Println("Строка:", v)
    case int:
        fmt.Println("Число:", v)
    default:
        fmt.Printf("Другой тип: %T\n", v)
    }
}
```

### Практические применения

1. **Обработка различных типов данных** в generic функциях
2. **Сериализация/десериализация** с разными типами
3. **Валидация данных** с различными форматами
4. **Паттерн Visitor** для обработки разных типов

### Дополнительные материалы

- [Go Tour: Type switches](https://tour.golang.org/methods/16)
- [Effective Go: Interface conversions and type assertions](https://golang.org/doc/effective_go.html#interface_conversions)
- [Go by Example: Type Switches](https://gobyexample.com/type-switches) 