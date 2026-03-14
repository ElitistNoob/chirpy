package chirps

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
)

func GetAllHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbChirps, err := appState.Queries.GetAllChirps(r.Context())
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}

		var chirps []chirp
		for _, c := range dbChirps {
			chirps = append(chirps, chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdateAt:  c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}

		internal.RespondWithJSON(w, http.StatusOK, chirps)
	}
}
