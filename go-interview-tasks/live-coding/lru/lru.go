// Пакет lru — кэш с вытеснением давно не используемых элементов (Least Recently Used).
//
// Идея: map даёт доступ к узлу за O(1), двусвязный список хранит порядок «свежести».
// Голова — недавно трогали, хвост — кандидат на удаление при переполнении.
package lru

// Cache — LRU-кэш фиксированной ёмкости.
type Cache struct {
	capacity int
	cache    map[int]*node // ключ → узел списка, O(1) поиск
	head     *node         // самый «свежий» элемент
	tail     *node         // самый «старый», выкидываем при переполнении
}

// node — элемент двусвязного списка.
type node struct {
	key, value int
	prev, next *node
}

// New создаёт кэш на capacity элементов.
func New(capacity int) *Cache {
	if capacity <= 0 {
		capacity = 1
	}
	return &Cache{
		capacity: capacity,
		cache:    make(map[int]*node, capacity),
	}
}

// Get возвращает значение по ключу. При попадании поднимает элемент в голову списка.
func (c *Cache) Get(key int) (int, bool) {
	n, ok := c.cache[key]
	if !ok {
		return 0, false
	}
	c.moveToHead(n) // прочитали — элемент стал недавно использованным
	return n.value, true
}

// Put добавляет или обновляет значение. При переполнении удаляет хвост списка.
func (c *Cache) Put(key, value int) {
	if n, ok := c.cache[key]; ok {
		n.value = value
		c.moveToHead(n) // ключ уже есть — обновляем и поднимаем в голову
		return
	}
	if len(c.cache) >= c.capacity {
		c.removeTail() // места нет — удаляем LRU-элемент
	}
	n := &node{key: key, value: value}
	c.cache[key] = n
	c.addToHead(n)
}

// addToHead вставляет узел в начало списка.
func (c *Cache) addToHead(n *node) {
	n.prev, n.next = nil, c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

// removeNode вырезает узел из списка (map не трогаем).
func (c *Cache) removeNode(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

// moveToHead переносит существующий узел в голову — O(1).
func (c *Cache) moveToHead(n *node) {
	c.removeNode(n)
	c.addToHead(n)
}

// removeTail удаляет самый старый элемент из кэша и списка.
func (c *Cache) removeTail() {
	if c.tail == nil {
		return
	}
	delete(c.cache, c.tail.key)
	c.removeNode(c.tail)
}
