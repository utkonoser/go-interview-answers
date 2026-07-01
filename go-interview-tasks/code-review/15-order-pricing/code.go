//go:build !solution

// Задача на code review (уровень: Lamoda).
// Расчёт суммы заказа в checkout — цены в каталоге e-commerce.
package orderpricing

// Line — позиция в корзине.
type Line struct {
	PriceRub float64 // цена в рублях, например 1499.99
	Qty      int
}

// TotalKopecks возвращает итог в копейках после скидки percent (0–100).
func TotalKopecks(lines []Line, discountPercent float64) int64 {
	var rub float64
	for _, l := range lines {
		rub += l.PriceRub * float64(l.Qty)
	}
	rub = rub * (100 - discountPercent) / 100
	return int64(rub * 100)
}
