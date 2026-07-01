//go:build solution

package waitgroupmisuse

import "sync"

// fix: wg.Add(1) внутри горутины — race с wg.Wait(): Wait может вернуться
// до того, как все Add выполнены → не все work() дождутся.
func Process(n int, work func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			work(id)
		}(i)
	}

	wg.Wait()
}
