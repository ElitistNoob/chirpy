package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const TokenTypeAccess TokenType = "chirpy-access"

func HashPassword(pw string) (string, error) {
	hash, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	return hash, err
}

func CheckPassword(pw, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(pw, hash)
	return match, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString(signingKey)
}

var parser = jwt.NewParser(
	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := parser.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("unexpected claims type")
	}

	if claims.Issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID in token subject: %w", err)
	}

	return userID, nil
}
