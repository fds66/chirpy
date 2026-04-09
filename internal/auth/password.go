package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
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

func GetBearerToken(headers http.Header) (string, error) {
	//fmt.Printf("received header %v\n", headers)

	//look for authorization header , Bearer TOKEN_STRING
	bearerToken := headers.Get("Authorization")
	//fmt.Printf("bearer token string %s\n", bearerToken)
	if bearerToken == "" {
		errorMessage := fmt.Errorf("No bearer token found ")
		fmt.Println("empty string from get")
		return "", errorMessage
	}
	if !strings.Contains(bearerToken, "Bearer") {
		errorMessage := fmt.Errorf("No bearer token found ")
		fmt.Println("no Bearer found in string")
		return "", errorMessage
	}
	bearerToken = strings.Replace(bearerToken, "Bearer", "", 1)
	bearerToken = strings.TrimSpace(bearerToken)
	//fmt.Printf("final bearerToken %s\n", bearerToken)
	return bearerToken, nil

}

/*
type Header map[string][]string
func (h Header) Get(key string) string
*/
func MakeRefreshToken() string {
	//func Read(b []byte) (n int, err error)
	key := make([]byte, 32)
	rand.Read(key)
	token := hex.EncodeToString(key)
	return token
}
