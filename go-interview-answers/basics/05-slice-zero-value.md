# 5. Какое у slice zero value? Какие операции над ним возможны?

## Ответ

**Zero value слайса — `nil`.** С nil слайсом можно выполнять определенные операции.

### Zero Value

```go
var slice []int // slice == nil
fmt.Printf("slice == nil: %t\n", slice == nil)
fmt.Printf("len(slice): %d\n", len(slice))
fmt.Printf("cap(slice): %d\n", cap(slice))
```

### Операции с nil слайсом

#### 1. Функции len() и cap()
```go
var slice []int
fmt.Printf("len(nil slice): %d\n", len(slice)) // 0
fmt.Printf("cap(nil slice): %d\n", cap(slice)) // 0
```

#### 2. Функция append()
```go
var slice []int
slice = append(slice, 1, 2, 3)
fmt.Printf("После append: %v\n", slice) // [1 2 3]
```

#### 3. Range цикл
```go
var slice []int
for i, v := range slice {
    fmt.Printf("Индекс: %d, Значение: %d\n", i, v)
}
// Цикл не выполнится, так как len(slice) == 0
```

#### 4. Срезы (slicing)
```go
var slice []int
// Эти операции не вызовут панику
fmt.Printf("slice[:]: %v\n", slice[:])     // []
fmt.Printf("slice[0:0]: %v\n", slice[0:0]) // []
```

### Полный список операций со слайсами

```go
func demonstrateSliceOperations() {
    var slice []int
    
    // 1. len() - длина слайса
    fmt.Printf("len: %d\n", len(slice))
    
    // 2. cap() - емкость слайса
    fmt.Printf("cap: %d\n", cap(slice))
    
    // 3. append() - добавление элементов
    slice = append(slice, 1, 2, 3)
    fmt.Printf("После append: %v\n", slice)
    
    // 4. copy() - копирование элементов
    dest := make([]int, 3)
    copied := copy(dest, slice)
    fmt.Printf("Скопировано элементов: %d\n", copied)
    
    // 5. Срезы (slicing)
    fmt.Printf("slice[1:2]: %v\n", slice[1:2])
    fmt.Printf("slice[:2]: %v\n", slice[:2])
    fmt.Printf("slice[1:]: %v\n", slice[1:])
    
    // 6. Range цикл
    for i, v := range slice {
        fmt.Printf("Индекс: %d, Значение: %d\n", i, v)
    }
    
    // 7. Индексация (только если len > 0)
    if len(slice) > 0 {
        fmt.Printf("Первый элемент: %d\n", slice[0])
    }
}
```

### Примеры с nil слайсами

```go
func nilSliceExamples() {
    var slice []int
    
    // append работает с nil слайсами
    result := append(slice, 1, 2, 3)
    fmt.Printf("append(nil, 1,2,3): %v\n", result)
    
    // copy работает с nil слайсами
    dest := make([]int, 3)
    copied := copy(dest, slice)
    fmt.Printf("copy(dest, nil): скопировано %d элементов\n", copied)
    
    // Срезы работают с nil слайсами
    fmt.Printf("nil[:]: %v\n", slice[:])
    fmt.Printf("nil[0:0]: %v\n", slice[0:0])
    
    // Range работает с nil слайсами (не выполняется)
    for i, v := range slice {
        fmt.Printf("Это не выполнится: %d, %d\n", i, v)
    }
    fmt.Println("Range с nil слайсом не выполняется")
}
```

### Сравнение с пустым слайсом

```go
func compareNilAndEmpty() {
    var nilSlice []int
    emptySlice := []int{}
    
    fmt.Printf("nilSlice == nil: %t\n", nilSlice == nil)
    fmt.Printf("emptySlice == nil: %t\n", emptySlice == nil)
    
    fmt.Printf("len(nilSlice): %d\n", len(nilSlice))
    fmt.Printf("len(emptySlice): %d\n", len(emptySlice))
    
    fmt.Printf("cap(nilSlice): %d\n", cap(nilSlice))
    fmt.Printf("cap(emptySlice): %d\n", cap(emptySlice))
    
    // Оба работают одинаково с append
    nilSlice = append(nilSlice, 1)
    emptySlice = append(emptySlice, 1)
    
    fmt.Printf("После append - nilSlice: %v\n", nilSlice)
    fmt.Printf("После append - emptySlice: %v\n", emptySlice)
}
```

### Практические примеры

#### Обработка опциональных данных
```go
func processOptionalData(data []int) {
    if data == nil {
        fmt.Println("Данные не предоставлены")
        return
    }
    
    if len(data) == 0 {
        fmt.Println("Данные пустые")
        return
    }
    
    fmt.Printf("Обрабатываем %d элементов\n", len(data))
}
```

#### Безопасное копирование
```go
func safeCopy(src []int) []int {
    if src == nil {
        return nil
    }
    
    dest := make([]int, len(src))
    copy(dest, src)
    return dest
}
```

#### Проверка на пустоту
```go
func isEmpty(slice []int) bool {
    return slice == nil || len(slice) == 0
}
```

### Важные моменты

1. **nil слайс функционально эквивалентен** пустому слайсу для большинства операций
2. **append работает с nil слайсами** и создает новый слайс
3. **len() и cap() возвращают 0** для nil слайса
4. **Range цикл не выполняется** для nil слайса
5. **Индексация вызовет панику** для nil слайса

### Дополнительные материалы

- [Go Tour: Slices](https://tour.golang.org/moretypes/8)
- [Effective Go: Slices](https://golang.org/doc/effective_go.html#slices)
- [Go by Example: Slices](https://gobyexample.com/slices) 