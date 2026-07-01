//go:build solution

package parindex

import (
	"strings"
	"sync"
)

func FindIndexes(docs []string, query string) []int {
	if query == "" {
		return nil
	}

	q := strings.ToLower(query)
	var (
		mu      sync.Mutex
		results []int
		wg      sync.WaitGroup
	)

	for i, doc := range docs {
		wg.Add(1)
		go func(idx int, text string) {
			defer wg.Done()
			if strings.Contains(strings.ToLower(text), q) {
				mu.Lock()
				results = append(results, idx)
				mu.Unlock()
			}
		}(i, doc)
	}

	wg.Wait()
	return results
}

// fix: захват i/doc из цикла; гонка на results; в MaxEven — гонка на max.
func MaxEven(nums []int) int {
	var (
		mu  sync.Mutex
		max int
		wg  sync.WaitGroup
	)

	for _, n := range nums {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			if v%2 == 0 {
				mu.Lock()
				if v > max {
					max = v
				}
				mu.Unlock()
			}
		}(n)
	}

	wg.Wait()
	return max
}
