//go:build !solution

// Задача на code review: потокобезопасный счётчик для метрик.
// Найдите проблемы перед мержем в production.
package racecounter

// Counter считает события из нескольких горутин.
type Counter struct {
	n int
}

func (c *Counter) Inc() {
	c.n++
}

func (c *Counter) Value() int {
	return c.n
}
