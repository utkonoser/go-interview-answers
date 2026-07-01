//go:build !solution

// Задача на code review (уровень: 2ГИС Application Review).
// REST-сервис: middleware считает запросы к API.
// На собесе 2ГИС просят «реализовать счётчик запросов к ручке средствами языка» — здесь типичные ошибки.
package requestmetrics

import "net/http"

// Metrics считает обработанные HTTP-запросы.
type Metrics struct {
	total int
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.total++
		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Total() int {
	return m.total
}

// RegisterRoutes — демо-хендлеры сервиса справочника POI (контекст 2ГИС).
func RegisterRoutes(mux *http.ServeMux, m *Metrics) {
	mux.Handle("/api/v1/places", m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})))
}
