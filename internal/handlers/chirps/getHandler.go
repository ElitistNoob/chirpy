package chirps

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/google/uuid"
)

func GetHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
			return
		}

		dbChirp, err := appState.Queries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp", err)
			return
		}

		internal.RespondWithJSON(w, http.StatusOK, chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdateAt:  dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}
}
