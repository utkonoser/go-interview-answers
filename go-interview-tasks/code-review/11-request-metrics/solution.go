//go:build solution

package requestmetrics

import (
	"net/http"
	"sync/atomic"
)

// fix: data race на total++ при конкурентных запросах.
type Metrics struct {
	total atomic.Int64
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.total.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Total() int64 {
	return m.total.Load()
}

func RegisterRoutes(mux *http.ServeMux, m *Metrics) {
	mux.Handle("/api/v1/places", m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})))
}
