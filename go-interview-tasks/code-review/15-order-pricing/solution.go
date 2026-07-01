//go:build solution

package orderpricing

import "math"

// Line — позиция в корзине, цена в копейках (как в платёжных системах).
type Line struct {
	PriceKopecks int64
	Qty          int
}

// fix: float64 для денег → ошибки округления и потеря копеек при int64().
// fix: скидку считать в integer: total * (100 - pct) / 100 с half-up при необходимости.
func TotalKopecks(lines []Line, discountPercent int64) int64 {
	var sum int64
	for _, l := range lines {
		sum += l.PriceKopecks * int64(l.Qty)
	}
	if discountPercent <= 0 {
		return sum
	}
	if discountPercent >= 100 {
		return 0
	}
	return (sum * (100 - discountPercent) + 50) / 100
}

// RubToKopecks — хелпер для тестов/миграции с API в рублях (округление half-up).
func RubToKopecks(rub float64) int64 {
	return int64(math.Round(rub * 100))
}
