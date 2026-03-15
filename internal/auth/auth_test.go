package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCreateHash(t *testing.T) {
	password := "secret-pa$$word"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash to be generated")
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("invalid hash format")
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		t.Fatalf("comparison failed: %s", err)
	}

	if !match {
		t.Fatalf("expected password to match hash")
	}
}

func TestCreateWrongHash(t *testing.T) {
	password := "secret-pa$$word"
	wrongPassword := "secret-password"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	match, err := argon2id.ComparePasswordAndHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("comparison failed: %v", err)
	}

	if match {
		t.Fatal("expected password and hash not to match")
	}
}

func TestValidateJWT(t *testing.T) {
	userID1 := uuid.New()
	secret1 := "correct-secret"
	secret2 := "wrong-secret"
	validToken, _ := MakeJWT(userID1, secret1, time.Minute*5)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		Subject:   userID1.String(),
	})
	expiredToken, _ := token.SignedString([]byte(secret1))

	tests := []struct {
		name      string
		token     string
		secret    string
		expectID  uuid.UUID
		expectErr bool
	}{
		{
			name:      "Valid Token",
			token:     validToken,
			secret:    secret1,
			expectID:  userID1,
			expectErr: false,
		},
		{
			name:      "Wrong Secret",
			token:     validToken,
			secret:    secret2,
			expectID:  uuid.UUID{},
			expectErr: true,
		},
		{
			name:      "Expired Token",
			token:     expiredToken,
			secret:    secret1,
			expectID:  userID1,
			expectErr: true,
		},
		{
			name:      "Malformed Token",
			token:     "not-a-token",
			secret:    secret1,
			expectID:  userID1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultingID, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateJWT() err = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && resultingID != tt.expectID {
				t.Errorf("ValidateJWT() expect = %v, got %v", resultingID, tt.expectID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		expectErr bool
	}{
		{"Valid Bearer token", "Bearer abc123", false},
		{"Missing Bearer", "abc123", true},
		{"Empty Header", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			header.Set("Authorization", tt.header)

			_, err := GetBearerToken(header)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetBearerToken() error = %v, expected: %v", err, tt.expectErr)
			}
		})
	}
}
