# 6. Как устроен тип map?

## Ответ

**Map в Go — это хеш-таблица**, которая обеспечивает амортизированную O(1) сложность для операций вставки, поиска и удаления.

### Базовая структура

Map состоит из:
1. **Массив buckets** (ведер) - основная структура данных
2. **Хеш-функция** для вычисления индекса
3. **Списки переполнения** для разрешения коллизий

### Устройство map

```go
// Упрощенная схема внутренней структуры map
type hmap struct {
    count     int    // количество элементов
    flags     uint8
    B         uint8  // log_2 количества buckets
    noverflow uint16 // количество overflow buckets
    hash0     uint32 // seed для хеш-функции
    
    buckets    unsafe.Pointer // массив из 2^B buckets
    oldbuckets unsafe.Pointer // предыдущие buckets при ресайзе
    nevacuate  uintptr        // количество эвакуированных old buckets
    
    extra *mapextra // дополнительные поля
}

type bmap struct {
    tophash [8]uint8 // верхние 8 бит хеша для каждого ключа
    // Далее идут ключи и значения (8 пар)
    // keys    [8]keytype
    // values  [8]valuetype
    // overflow *bmap // указатель на следующий bucket при переполнении
}
```

### Хеш-функция

```go
// Go использует быструю хеш-функцию
// Для строк: алгоритм, оптимизированный для ASCII
// Для чисел: простая, но эффективная функция
// Для структур: комбинирует хеши полей

func demonstrateHash() {
    // Хеш для строк
    str := "hello"
    fmt.Printf("Хеш строки '%s': %d\n", str, hashString(str))
    
    // Хеш для чисел
    num := 42
    fmt.Printf("Хеш числа %d: %d\n", num, hashInt(num))
}

func hashString(s string) uint32 {
    // Упрощенная демонстрация
    h := uint32(0)
    for i := 0; i < len(s); i++ {
        h = 31*h + uint32(s[i])
    }
    return h
}

func hashInt(i int) uint32 {
    // Упрощенная демонстрация
    return uint32(i)
}
```

### Bucket (ведро)

```go
// Каждый bucket может содержать до 8 пар ключ-значение
// Если bucket переполняется, создается overflow bucket

func demonstrateBucket() {
    m := make(map[string]int)
    
    // Добавляем элементы
    m["a"] = 1
    m["b"] = 2
    m["c"] = 3
    m["d"] = 4
    m["e"] = 5
    m["f"] = 6
    m["g"] = 7
    m["h"] = 8
    m["i"] = 9 // Этот элемент попадет в overflow bucket
    
    fmt.Printf("Map содержит %d элементов\n", len(m))
}
```

### Операции с map

#### Создание
```go
// Пустой map
m1 := make(map[string]int)

// Map с начальной capacity
m2 := make(map[string]int, 100)

// Map с литералом
m3 := map[string]int{
    "a": 1,
    "b": 2,
}
```

#### Вставка и обновление
```go
func mapOperations() {
    m := make(map[string]int)
    
    // Вставка
    m["key1"] = 1
    
    // Обновление
    m["key1"] = 2
    
    // Проверка существования
    value, exists := m["key1"]
    if exists {
        fmt.Printf("Значение: %d\n", value)
    }
}
```

#### Удаление
```go
func mapDelete() {
    m := map[string]int{"a": 1, "b": 2}
    
    // Удаление элемента
    delete(m, "a")
    
    // Проверка после удаления
    if _, exists := m["a"]; !exists {
        fmt.Println("Элемент 'a' удален")
    }
}
```

### Коллизии и разрешение

```go
// Когда хеши разных ключей указывают на один bucket
// Go использует метод цепочек (chaining)

func demonstrateCollisions() {
    m := make(map[string]int)
    
    // Добавляем элементы, которые могут попасть в один bucket
    for i := 0; i < 20; i++ {
        key := fmt.Sprintf("key%d", i)
        m[key] = i
    }
    
    fmt.Printf("Map содержит %d элементов\n", len(m))
}
```

### Ресайз map

```go
// Map автоматически увеличивается при достижении load factor
// Load factor = количество элементов / количество buckets

func demonstrateResize() {
    m := make(map[int]int)
    
    fmt.Printf("Начальный размер: %d\n", len(m))
    
    // Добавляем элементы до ресайза
    for i := 0; i < 100; i++ {
        m[i] = i * 2
        if i%10 == 0 {
            fmt.Printf("После %d элементов: len=%d\n", i, len(m))
        }
    }
}
```

### Производительность

```go
func benchmarkMapOperations() {
    m := make(map[int]int, 1000)
    
    // Вставка
    start := time.Now()
    for i := 0; i < 1000; i++ {
        m[i] = i * 2
    }
    insertTime := time.Since(start)
    
    // Поиск
    start = time.Now()
    for i := 0; i < 1000; i++ {
        _ = m[i]
    }
    lookupTime := time.Since(start)
    
    fmt.Printf("Вставка 1000 элементов: %v\n", insertTime)
    fmt.Printf("Поиск 1000 элементов: %v\n", lookupTime)
}
```

### Особенности map в Go

#### 1. Неупорядоченность
```go
func demonstrateUnordered() {
    m := map[string]int{
        "a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
    }
    
    fmt.Println("Порядок итерации может меняться:")
    for k, v := range m {
        fmt.Printf("%s: %d\n", k, v)
    }
}
```

#### 2. Нельзя получить адрес элемента
```go
func demonstrateNoAddress() {
    m := map[string]int{"a": 1}
    
    // Это не скомпилируется:
    // ptr := &m["a"] // ошибка компиляции
    
    // Нужно использовать промежуточную переменную
    value := m["a"]
    ptr := &value
    fmt.Printf("Значение через указатель: %d\n", *ptr)
}
```

#### 3. Сравнение map
```go
func demonstrateMapComparison() {
    m1 := map[string]int{"a": 1}
    m2 := map[string]int{"a": 1}
    
    // Можно сравнивать только с nil
    fmt.Printf("m1 == nil: %t\n", m1 == nil)
    
    // Для сравнения содержимого нужно писать свою функцию
    fmt.Printf("Maps равны: %t\n", mapsEqual(m1, m2))
}

func mapsEqual(m1, m2 map[string]int) bool {
    if len(m1) != len(m2) {
        return false
    }
    for k, v := range m1 {
        if m2[k] != v {
            return false
        }
    }
    return true
}
```

### Практические советы

1. **Используйте make с capacity** для больших map
2. **Проверяйте существование** ключа при чтении
3. **Помните о неупорядоченности** при итерации
4. **Используйте sync.Map** для конкурентного доступа

### Дополнительные материалы

- [Go Tour: Maps](https://tour.golang.org/moretypes/19)
- [Effective Go: Maps](https://golang.org/doc/effective_go.html#maps)
- [Go by Example: Maps](https://gobyexample.com/maps)
- [Go Maps in Action](https://blog.golang.org/maps) 