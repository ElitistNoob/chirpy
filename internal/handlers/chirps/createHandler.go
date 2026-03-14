package chirps

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	db "github.com/ElitistNoob/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	UserID uuid.UUID `json:"user_id"`
	Body   string    `json:"body"`
}

func CreateHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params parameters
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Error decoding body response", err)
			return
		}

		cleaned, err := validateChirp(params.Body)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "chirp is invalid", err)
			return
		}

		c, err := appState.Queries.CreateChirp(r.Context(), db.CreateChirpParams{
			UserID: params.UserID,
			Body:   cleaned,
		})
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "Failed creating chirp", err)
			return
		}

		internal.RespondWithJSON(w, http.StatusCreated, chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdateAt:  c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
}

func validateChirp(body string) (string, error) {
	const MAX_LEN = 140

	if len(body) > MAX_LEN {
		err := errors.New("Chirp cannot be longer than 140 characters")
		return "", err
	}

	return sanitizeProfanity(body), nil
}

func sanitizeProfanity(msg string) string {
	profanity := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(msg, " ")

	for i, word := range words {
		l_word := strings.ToLower(word)
		if _, ok := profanity[l_word]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
