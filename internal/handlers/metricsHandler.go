package handlers

import (
	"fmt"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal/app"
)

func MiddlewareMetricsInc(app *app.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
	}
}

func MetricsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		value := app.FileserverHits.Load()
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
}
