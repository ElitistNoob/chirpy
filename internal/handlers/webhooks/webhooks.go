package webhooks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/ElitistNoob/chirpy/internal/auth"
	"github.com/google/uuid"
)

func WebhookHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const event = "user.upgraded"
		type parameters struct {
			Event string `json:"event"`
			Data  struct {
				UserID uuid.UUID `json:"user_id"`
			}
		}

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			internal.RespondWithError(w, http.StatusUnauthorized, "couldn't find apikey", err)
			return
		}

		if apiKey != appState.PolkaKey {
			internal.RespondWithError(w, http.StatusUnauthorized, "apikey is invalid", err)
			return
		}

		var params parameters
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&params)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "Couldn't decode params", err)
			return
		}

		if params.Event != event {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		_, err = appState.Queries.UpdateChirpyRedStatus(r.Context(), params.Data.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				internal.RespondWithError(w, http.StatusNotFound, "User not found", err)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
