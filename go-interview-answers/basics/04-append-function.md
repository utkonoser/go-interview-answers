# 4. Как работает append?

## Ответ

**`append`** — это встроенная функция в Go для добавления элементов в слайс. Она может изменять capacity слайса при необходимости.

### Базовый синтаксис

```go
slice = append(slice, element1, element2, ...)
```

### Как работает append

1. **Проверяет capacity**: Если в слайсе достаточно места, просто добавляет элементы
2. **Увеличивает capacity**: Если места недостаточно, создает новый массив с большим capacity
3. **Копирует данные**: Копирует все существующие элементы в новый массив
4. **Добавляет новые элементы**: Добавляет новые элементы в конец
5. **Возвращает новый слайс**: Возвращает слайс, указывающий на новый массив

### Примеры

#### Простое добавление
```go
package main

import "fmt"

func main() {
    // Создаем слайс с capacity 3
    slice := make([]int, 0, 3)
    fmt.Printf("До append: len=%d, cap=%d, %v\n", len(slice), cap(slice), slice)
    
    // Добавляем элементы
    slice = append(slice, 1)
    fmt.Printf("После append(1): len=%d, cap=%d, %v\n", len(slice), cap(slice), slice)
    
    slice = append(slice, 2, 3)
    fmt.Printf("После append(2,3): len=%d, cap=%d, %v\n", len(slice), cap(slice), slice)
    
    // Добавляем еще один элемент - capacity увеличится
    slice = append(slice, 4)
    fmt.Printf("После append(4): len=%d, cap=%d, %v\n", len(slice), cap(slice), slice)
}
```

#### Стратегия роста capacity

```go
func demonstrateGrowth() {
    var slice []int
    
    for i := 0; i < 20; i++ {
        slice = append(slice, i)
        fmt.Printf("len=%d, cap=%d\n", len(slice), cap(slice))
    }
}
```

### Стратегия роста capacity

Go использует следующую стратегию:
- **До 1024 элементов**: capacity удваивается
- **После 1024 элементов**: capacity увеличивается на 25%

```go
func showGrowthStrategy() {
    var slice []int
    
    fmt.Println("len\tcap\tgrowth")
    fmt.Println("---\t---\t------")
    
    for i := 0; i < 2000; i++ {
        oldCap := cap(slice)
        slice = append(slice, i)
        newCap := cap(slice)
        
        if newCap != oldCap {
            growth := float64(newCap) / float64(oldCap)
            fmt.Printf("%d\t%d\t%.2fx\n", len(slice), newCap, growth)
        }
    }
}
```

### Работа с указателями

```go
func demonstratePointerBehavior() {
    // Создаем слайс
    original := []int{1, 2, 3}
    fmt.Printf("Original: %v, ptr: %p\n", original, &original[0])
    
    // Добавляем элемент
    modified := append(original, 4)
    fmt.Printf("Modified: %v, ptr: %p\n", modified, &modified[0])
    
    // Проверяем, изменился ли указатель
    if &original[0] == &modified[0] {
        fmt.Println("Указатель не изменился - capacity было достаточно")
    } else {
        fmt.Println("Указатель изменился - создан новый массив")
    }
}
```

### Множественное добавление

```go
func multipleAppend() {
    slice1 := []int{1, 2, 3}
    slice2 := []int{4, 5, 6}
    
    // Добавляем элементы из другого слайса
    result := append(slice1, slice2...)
    fmt.Printf("Результат: %v\n", result)
}
```

### Работа с nil слайсами

```go
func nilSliceAppend() {
    var slice []int // nil слайс
    
    fmt.Printf("Nil slice: %v, len=%d, cap=%d\n", slice, len(slice), cap(slice))
    
    // append работает с nil слайсами
    slice = append(slice, 1, 2, 3)
    fmt.Printf("После append: %v, len=%d, cap=%d\n", slice, len(slice), cap(slice))
}
```

### Эффективность append

```go
func efficientAppend() {
    // Неэффективно - много перераспределений
    var slice1 []int
    for i := 0; i < 1000; i++ {
        slice1 = append(slice1, i)
    }
    
    // Эффективно - предварительно выделяем память
    slice2 := make([]int, 0, 1000)
    for i := 0; i < 1000; i++ {
        slice2 = append(slice2, i)
    }
    
    fmt.Printf("Неэффективный: %d перераспределений\n", countReallocations(slice1))
    fmt.Printf("Эффективный: %d перераспределений\n", countReallocations(slice2))
}

func countReallocations(slice []int) int {
    // Упрощенная логика подсчета перераспределений
    return 0 // Реальная реализация сложнее
}
```

### Работа с capacity

```go
func capacityExamples() {
    // Создаем слайс с определенной capacity
    slice := make([]int, 0, 5)
    fmt.Printf("Начальный: len=%d, cap=%d\n", len(slice), cap(slice))
    
    // Добавляем элементы до capacity
    for i := 0; i < 5; i++ {
        slice = append(slice, i)
        fmt.Printf("После %d: len=%d, cap=%d\n", i, len(slice), cap(slice))
    }
    
    // Добавляем еще один - capacity увеличится
    slice = append(slice, 5)
    fmt.Printf("После превышения capacity: len=%d, cap=%d\n", len(slice), cap(slice))
}
```

### Практические советы

1. **Предварительно выделяйте память** для больших слайсов
2. **Используйте `make`** с указанием capacity
3. **Помните о копировании** при превышении capacity
4. **Используйте `copy`** для копирования слайсов

### Дополнительные материалы

- [Go Tour: Slices](https://tour.golang.org/moretypes/8)
- [Effective Go: Slices](https://golang.org/doc/effective_go.html#slices)
- [Go by Example: Slices](https://gobyexample.com/slices) 