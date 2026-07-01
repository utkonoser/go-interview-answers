//go:build !solution

// Задача на code review: параллельная обработка элементов.
package waitgroupmisuse

import "sync"

// Process запускает n горутин и ждёт их завершения.
func Process(n int, work func(int)) {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		go func(id int) {
			wg.Add(1)
			defer wg.Done()
			work(id)
		}(i)
	}

	wg.Wait()
}
