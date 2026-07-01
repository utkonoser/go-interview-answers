//go:build solution

package goroutineleak

// fix: канал не закрыт → range блокируется навсегда → утечка горутины-отправителя.
func DoubleAll(items []int) []int {
	ch := make(chan int, len(items))

	go func() {
		defer close(ch)
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
