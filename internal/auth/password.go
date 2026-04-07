package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
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
