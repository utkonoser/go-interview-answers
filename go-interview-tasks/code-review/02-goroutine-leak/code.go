//go:build !solution

// Задача на code review: удвоить числа из слайса в фоне.
package goroutineleak

// DoubleAll запускает воркер и собирает результаты.
func DoubleAll(items []int) []int {
	ch := make(chan int)

	go func() {
		for _, item := range items {
			ch <- item * 2
		}
	}()

	out := make([]int, 0, len(items))
	for v := range ch {
		out = append(out, v)
	}
	return out
}
