package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
	"github.com/ElitistNoob/chirpy/internal/database"
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
		if err != nil || !match {
			internal.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		token, err := auth.MakeJWT(dbUser.ID, appState.Secret, time.Hour)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "Couldn't create new token", err)
			return
		}

		refreshToken, err := appState.Queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     auth.MakeRefreshToken(),
			UserID:    dbUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		})
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Couldn't create new refreshToken", err)
			return
		}

		internal.RespondWithJSON(w, http.StatusOK, User{
			ID:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			IsChirpyRed:  dbUser.IsChirpyRed,
			Token:        token,
			RefreshToken: refreshToken.Token,
		})
	}
}
