package metrics

import (
	"net/http"
)

func (m *Metrics) MiddlewareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
