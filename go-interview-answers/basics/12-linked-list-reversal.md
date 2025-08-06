# 12. Сколько времени в минутах займет у вас написание процедуры обращения односвязного списка?

## Ответ

**Написание процедуры разворота односвязного списка занимает 2-5 минут для опытного разработчика.** Это классическая задача, которая проверяет понимание структур данных и алгоритмов.

### Определение односвязного списка

```go
package main

import "fmt"

type Node struct {
    Value int
    Next  *Node
}

type LinkedList struct {
    Head *Node
}
```

### Итеративное решение (рекомендуемое)

```go
func (list *LinkedList) Reverse() {
    var prev *Node = nil
    current := list.Head
    
    for current != nil {
        // Сохраняем следующий узел
        next := current.Next
        
        // Меняем направление указателя
        current.Next = prev
        
        // Переходим к следующему узлу
        prev = current
        current = next
    }
    
    // Обновляем голову списка
    list.Head = prev
}
```

### Рекурсивное решение

```go
func (list *LinkedList) ReverseRecursive() {
    list.Head = reverseRecursiveHelper(list.Head, nil)
}

func reverseRecursiveHelper(current, prev *Node) *Node {
    if current == nil {
        return prev
    }
    
    // Сохраняем следующий узел
    next := current.Next
    
    // Меняем направление указателя
    current.Next = prev
    
    // Рекурсивно обрабатываем остальную часть списка
    return reverseRecursiveHelper(next, current)
}
```

### Полный пример с тестированием

```go
package main

import "fmt"

type Node struct {
    Value int
    Next  *Node
}

type LinkedList struct {
    Head *Node
}

// Создание нового узла
func NewNode(value int) *Node {
    return &Node{Value: value, Next: nil}
}

// Добавление элемента в конец списка
func (list *LinkedList) Append(value int) {
    newNode := NewNode(value)
    
    if list.Head == nil {
        list.Head = newNode
        return
    }
    
    current := list.Head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
}

// Итеративный разворот
func (list *LinkedList) Reverse() {
    var prev *Node = nil
    current := list.Head
    
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    
    list.Head = prev
}

// Рекурсивный разворот
func (list *LinkedList) ReverseRecursive() {
    list.Head = reverseRecursiveHelper(list.Head, nil)
}

func reverseRecursiveHelper(current, prev *Node) *Node {
    if current == nil {
        return prev
    }
    
    next := current.Next
    current.Next = prev
    return reverseRecursiveHelper(next, current)
}

// Вывод списка
func (list *LinkedList) Print() {
    current := list.Head
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

// Получение длины списка
func (list *LinkedList) Length() int {
    count := 0
    current := list.Head
    for current != nil {
        count++
        current = current.Next
    }
    return count
}
```

### Демонстрация работы

```go
func main() {
    // Создаем список: 1 -> 2 -> 3 -> 4 -> 5
    list := &LinkedList{}
    list.Append(1)
    list.Append(2)
    list.Append(3)
    list.Append(4)
    list.Append(5)
    
    fmt.Println("Исходный список:")
    list.Print()
    
    // Разворачиваем список
    list.Reverse()
    
    fmt.Println("Развернутый список:")
    list.Print()
    
    // Разворачиваем обратно
    list.Reverse()
    
    fmt.Println("Восстановленный список:")
    list.Print()
}
```

### Тестирование

```go
import "testing"

func TestLinkedListReverse(t *testing.T) {
    // Создаем тестовый список
    list := &LinkedList{}
    values := []int{1, 2, 3, 4, 5}
    
    for _, value := range values {
        list.Append(value)
    }
    
    // Проверяем исходный порядок
    current := list.Head
    for i, expected := range values {
        if current == nil {
            t.Errorf("Ожидался узел с индексом %d", i)
            return
        }
        if current.Value != expected {
            t.Errorf("Ожидалось %d, получено %d", expected, current.Value)
        }
        current = current.Next
    }
    
    // Разворачиваем список
    list.Reverse()
    
    // Проверяем развернутый порядок
    current = list.Head
    for i := len(values) - 1; i >= 0; i-- {
        if current == nil {
            t.Errorf("Ожидался узел с индексом %d", i)
            return
        }
        if current.Value != values[i] {
            t.Errorf("Ожидалось %d, получено %d", values[i], current.Value)
        }
        current = current.Next
    }
}
```

### Сложность алгоритма

```go
func analyzeComplexity() {
    fmt.Println("Анализ сложности разворота односвязного списка:")
    fmt.Println("Временная сложность: O(n)")
    fmt.Println("Пространственная сложность: O(1) для итеративного решения")
    fmt.Println("Пространственная сложность: O(n) для рекурсивного решения (стек вызовов)")
}
```

### Варианты реализации

#### 1. Разворот части списка
```go
func (list *LinkedList) ReverseBetween(start, end int) {
    if start >= end || list.Head == nil {
        return
    }
    
    // Находим начальный узел
    dummy := &Node{Next: list.Head}
    prev := dummy
    
    for i := 0; i < start; i++ {
        prev = prev.Next
    }
    
    current := prev.Next
    
    // Разворачиваем подсписок
    for i := 0; i < end-start; i++ {
        next := current.Next
        current.Next = next.Next
        next.Next = prev.Next
        prev.Next = next
    }
    
    list.Head = dummy.Next
}
```

#### 2. Разворот по группам
```go
func (list *LinkedList) ReverseInGroups(groupSize int) {
    if groupSize <= 1 {
        return
    }
    
    list.Head = reverseInGroupsHelper(list.Head, groupSize)
}

func reverseInGroupsHelper(head *Node, groupSize int) *Node {
    if head == nil {
        return nil
    }
    
    // Проверяем, достаточно ли узлов для группы
    current := head
    count := 0
    for current != nil && count < groupSize {
        current = current.Next
        count++
    }
    
    if count < groupSize {
        return head // Недостаточно узлов для разворота
    }
    
    // Разворачиваем текущую группу
    var prev *Node = nil
    current = head
    for i := 0; i < groupSize; i++ {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    
    // Рекурсивно обрабатываем следующую группу
    head.Next = reverseInGroupsHelper(current, groupSize)
    
    return prev
}
```

### Оптимизации

#### 1. Проверка на пустой список
```go
func (list *LinkedList) ReverseOptimized() {
    if list.Head == nil || list.Head.Next == nil {
        return // Пустой список или один элемент
    }
    
    var prev *Node = nil
    current := list.Head
    
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    
    list.Head = prev
}
```

#### 2. Разворот с сохранением длины
```go
func (list *LinkedList) ReverseWithLength() int {
    var prev *Node = nil
    current := list.Head
    length := 0
    
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
        length++
    }
    
    list.Head = prev
    return length
}
```

### Лучшие практики

```go
// 1. Всегда проверяйте граничные случаи
func (list *LinkedList) ReverseSafe() {
    if list.Head == nil {
        return
    }
    
    var prev *Node = nil
    current := list.Head
    
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    
    list.Head = prev
}

// 2. Используйте понятные имена переменных
func (list *LinkedList) ReverseClear() {
    var previousNode *Node = nil
    currentNode := list.Head
    
    for currentNode != nil {
        nextNode := currentNode.Next
        currentNode.Next = previousNode
        previousNode = currentNode
        currentNode = nextNode
    }
    
    list.Head = previousNode
}

// 3. Добавляйте комментарии для сложной логики
func (list *LinkedList) ReverseCommented() {
    // Алгоритм разворота: меняем направление всех указателей
    var prev *Node = nil
    current := list.Head
    
    for current != nil {
        // Сохраняем ссылку на следующий узел
        next := current.Next
        
        // Меняем направление указателя текущего узла
        current.Next = prev
        
        // Переходим к следующему узлу
        prev = current
        current = next
    }
    
    // Обновляем голову списка
    list.Head = prev
}
```

### Дополнительные материалы

- [Linked List Data Structure](https://en.wikipedia.org/wiki/Linked_list)
- [Go Data Structures](https://golang.org/doc/effective_go.html#data)
- [Algorithm Complexity](https://en.wikipedia.org/wiki/Time_complexity) 