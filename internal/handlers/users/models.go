package users

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	HashPassword string    `json:"hash_password"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}
