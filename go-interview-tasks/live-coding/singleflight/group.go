// Пакет singleflight — один in-flight запрос на ключ (защита от thundering herd).
//
// Когда кеш протух, 1000 горутин не должны одновременно идти в БД — только одна,
// остальные ждут её результат. Актуально для WB/Lamoda/Ozon.
package singleflight

import (
	"golang.org/x/sync/singleflight"
)

// Group дедуплицирует вызовы Load по ключу.
type Group struct {
	g singleflight.Group
}

// Load выполняет fn для key; параллельные вызовы с тем же key ждут один результат.
func (g *Group) Load(key string, fn func() (any, error)) (any, error) {
	v, err, _ := g.g.Do(key, func() (any, error) {
		return fn()
	})
	return v, err
}
