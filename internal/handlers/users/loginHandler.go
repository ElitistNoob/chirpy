package users

import (
	"encoding/json"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
)

func LoginHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode response body", err)
			return
		}

		dbUser, err := appState.Queries.GetUserByEmail(r.Context(), user.Email)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		match, err := auth.CheckPassword(user.Password, dbUser.HashedPassword)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Failed to authenticate user", err)
			return
		}

		if !match {
			internal.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		internal.RespondWithJSON(w, http.StatusOK, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	}
}
