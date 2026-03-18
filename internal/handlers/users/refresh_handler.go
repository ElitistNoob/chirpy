package users

import (
	"net/http"
	"time"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
)

func RefreshHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "couldn't get refresh token", err)
			return
		}

		user, err := appState.Queries.GetUserFromRefreshToken(r.Context(), refreshToken)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "couldn't get user from refresh token", err)
			return
		}

		accessToken, err := auth.MakeJWT(
			user.ID,
			appState.Secret,
			time.Hour,
		)
		if err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
			return
		}

		type Res struct {
			Token string `json:"token"`
		}

		internal.RespondWithJSON(w, http.StatusOK, Res{
			Token: accessToken,
		})
	}
}

func RevokeHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "couldn't get refresh token", err)
			return
		}

		if err := appState.Queries.RevokeResfreshToken(r.Context(), refreshToken); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "couldn't revoke token", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
