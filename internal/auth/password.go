package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {

	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		errorMessage := fmt.Errorf("Error hashing password %w", err)
		return "", errorMessage
	}
	return hashedPassword, nil

}

func CheckPasswordHash(password, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		errorMessage := fmt.Errorf("Error comparing password to hash %w", err)
		return false, errorMessage
	}
	return result, nil

}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	UTCTimeNow := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(UTCTimeNow),
		ExpiresAt: jwt.NewNumericDate(UTCTimeNow.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingSecret := []byte(tokenSecret)
	tokenString, err := token.SignedString(signingSecret)
	if err != nil {
		errorMessage := fmt.Errorf("Error creating token string, %w", err)
		return "", errorMessage
	}
	return tokenString, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	recClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &recClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		errorMessage := fmt.Errorf("Error creating token string %w", err)
		return uuid.UUID{}, errorMessage
	}
	user, err := token.Claims.GetSubject()
	if err != nil {
		errorMessage := fmt.Errorf("Error getting subject string %w", err)
		return uuid.UUID{}, errorMessage
	}
	userID, err := uuid.Parse(user)
	if err != nil {
		errorMessage := fmt.Errorf("Error converting subject string to uuid %w", err)
		return uuid.UUID{}, errorMessage
	}
	return userID, nil
}
