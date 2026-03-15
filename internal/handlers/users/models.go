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
	HashPassword string    `json:"hash_password"`
}
