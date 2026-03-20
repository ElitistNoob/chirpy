package users

import (
	"encoding/json"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
	"github.com/ElitistNoob/chirpy/internal/database"
)

func UpdateHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			User
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
			return
		}

		userID, err := auth.ValidateJWT(token, appState.Secret)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		var params UserRequest
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&params); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
			return
		}

		hashed_password, err := auth.HashPassword(params.Password)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
			return
		}

		updatedUser, err := appState.Queries.UpdateUserEmailAndPassword(
			r.Context(),
			database.UpdateUserEmailAndPasswordParams{
				Email:          params.Email,
				HashedPassword: hashed_password,
				ID:             userID,
			})
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "user couldn't be udpated", err)
			return
		}

		resp := response{
			User{
				ID:          updatedUser.ID,
				CreatedAt:   updatedUser.CreatedAt,
				UpdatedAt:   updatedUser.UpdatedAt,
				Email:       updatedUser.Email,
				IsChirpyRed: updatedUser.IsChirpyRed,
			},
		}

		internal.RespondWithJSON(w, http.StatusOK, resp)
	}
}
