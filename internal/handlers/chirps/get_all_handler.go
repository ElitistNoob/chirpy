package chirps

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/database"
	"github.com/google/uuid"
)

func GetAllHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID, err := getAuthorIdFromRequest(r)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "invalid author_id", err)
			return
		}

		isSortDesc := getSortOrderFromRequest(r)

		var dbChirps []database.Chirp
		request := database.GetChirpsByAuthorParams{
			UserID:      authorID,
			Isdescorder: isSortDesc,
		}

		if authorID != uuid.Nil {
			dbChirps, err = appState.Queries.GetChirpsByAuthor(r.Context(), request)
		} else {
			dbChirps, err = appState.Queries.GetAllChirps(r.Context(), isSortDesc)
		}

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
