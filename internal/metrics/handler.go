package metrics

import (
	"fmt"
	"net/http"
)

func (m *Metrics) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	value := m.fileserverHits.Load()
	htmlTemplate := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, value)
	w.Write([]byte(htmlTemplate))
}

func (m *Metrics) ResetHandler(w http.ResponseWriter, r *http.Request) {
	m.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
