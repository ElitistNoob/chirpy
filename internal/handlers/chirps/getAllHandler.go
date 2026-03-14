package chirps

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
)

func GetAllHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := appState.Queries.GetAllChirps(r.Context())
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "Error returning all chirps", err)
			return
		}

		var chirpsList []chirpModel
		for _, c := range chirps {
			chirpsList = append(chirpsList, chirpModel{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdateAt:  c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}

		internal.RespondWithJSON(w, http.StatusOK, chirpsList)
	}
}
