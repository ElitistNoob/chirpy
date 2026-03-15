package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
)

const (
	MAX_TOKEN_VALIDITY = 60 * 60
)

func LoginHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode response body", err)
			return
		}

		expiresIn := MAX_TOKEN_VALIDITY

		value := r.Header.Get("ExpiresInSeconds")
		if value != "" {
			if user.ExpiresInSeconds > 0 && user.ExpiresInSeconds < MAX_TOKEN_VALIDITY {
				expiresIn = user.ExpiresInSeconds
			}
		}

		dbUser, err := appState.Queries.GetUserByEmail(r.Context(), user.Email)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		match, err := auth.CheckPassword(user.Password, dbUser.HashedPassword)
		if err != nil || !match {
			internal.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		token, err := auth.MakeJWT(dbUser.ID, appState.Secret, time.Second*time.Duration(expiresIn))
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Couldn't create new token", err)
			return

		}

		internal.RespondWithJSON(w, http.StatusOK, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
			Token:     token,
		})
	}
}
