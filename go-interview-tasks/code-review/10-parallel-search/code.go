//go:build !solution

// Задача на code review (уровень: Wildberries/Ozon Middle).
// Параллельный поиск документов, содержащих подстроку.
// Формат «найди баги в коде коллеги» — см. Habr: реальные задачи WB/VK.
package parindex

import (
	"strings"
	"sync"
)

// FindIndexes возвращает индексы документов, где есть query (без учёта регистра).
func FindIndexes(docs []string, query string) []int {
	if query == "" {
		return nil
	}

	q := strings.ToLower(query)
	results := make([]int, 0)
	var wg sync.WaitGroup

	for i, doc := range docs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if strings.Contains(strings.ToLower(doc), q) {
				results = append(results, i)
			}
		}()
	}

	wg.Wait()
	return results
}

// MaxEven возвращает максимальное чётное число в nums (0 если нет).
func MaxEven(nums []int) int {
	max := 0
	var wg sync.WaitGroup

	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if n%2 == 0 && n > max {
				max = n
			}
		}()
	}

	wg.Wait()
	return max
}
