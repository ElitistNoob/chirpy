package users

import (
	"encoding/json"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
	"github.com/ElitistNoob/chirpy/internal/database"
)

func CreateHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
			return
		}

		user, err := appState.Queries.CreateUser(r.Context(), database.CreateUserParams{
			HashedPassword: hash,
			Email:          req.Email,
		})
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "bad request", err)
			return
		}

		res := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		}

		internal.RespondWithJSON(w, http.StatusCreated, res)
	}
}
