package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ElitistNoob/chirpy/internal"
	"github.com/ElitistNoob/chirpy/internal/app"
	"github.com/google/uuid"
)

type UserRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func CreateUserHandler(appState *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UserRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			internal.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		user, err := appState.Queries.CreateUser(r.Context(), req.Email)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, "bad request", err)
			return
		}

		res := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		internal.RespondWithJSON(w, http.StatusCreated, res)
	}
}
