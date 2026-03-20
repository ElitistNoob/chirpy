package chirps

import (
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
	"github.com/google/uuid"
)

func DeleteHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
			return
		}

		userID, err := auth.ValidateJWT(token, appState.Secret)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "invalid JWT", err)
			return
		}

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "Missing chirp Id", err)
			return
		}

		chirp, err := appState.Queries.GetChirpByID(r.Context(), chirpID)
		if err != nil {
			internal.RespondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}

		if chirp.UserID != userID {
			internal.RespondWithError(w, http.StatusForbidden, "insufficient permission to delete chirp", err)
			return
		}

		err = appState.Queries.DeleteChirp(r.Context(), chirpID)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Chirp couldn't be deleted", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
