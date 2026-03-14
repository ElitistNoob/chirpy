package handlers

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal/app"
)

func ResetHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appState.Platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Reset only allowed in dev environment"))
			return
		}

		appState.FileserverHits.Store(0)
		if err := appState.Queries.DeleteUsers(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to reset the database: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0 and database reset to initial state"))
	}
}
