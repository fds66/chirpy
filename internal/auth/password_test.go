package auth

import (
	//"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		input    string
		expected error
	}{
		{
			input:    "george",
			expected: nil,
		},
		{
			input:    "77gyskjhdyr",
			expected: nil,
		},

		// add more cases here
	}
	for _, c := range cases {
		_, err := HashPassword(c.input)

		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if err != c.expected {
			t.Errorf("error from hash function %sv", err)
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		input      string
		expected   bool
		decodeHash string
	}{
		{
			input:      "george",
			expected:   true,
			decodeHash: "",
		},
		{
			input:      "george",
			expected:   false,
			decodeHash: "lkjldjgekdfgj",
		},

		// add more cases here
	}
	for _, c := range cases {
		hashed, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("error from hash function %sv", err)
		}
		if c.decodeHash != "" {
			hashed = c.decodeHash
		}
		result, err := CheckPasswordHash(c.input, hashed)
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if result != c.expected {
			t.Errorf("password check gave inaccurate answer, actual %t, expected %t", result, c.expected)
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

	}
}

/*
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error)
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error)
*/
func TestValidateJWT(t *testing.T) {
	id1, _ := uuid.NewRandom()
	id2, _ := uuid.NewRandom()
	secret := "kkjsldkjflk"
	wrongSecret := "jjgjjgjgj"
	duration, _ := time.ParseDuration("1h")
	expired, _ := time.ParseDuration("-1h")

	cases := []struct {
		userID       uuid.UUID
		tokenSecret  string
		decodeSecret string
		expiresIn    time.Duration
		expected     uuid.UUID
	}{
		{
			userID:       id1,
			tokenSecret:  secret,
			decodeSecret: secret,
			expiresIn:    duration,
			expected:     id1,
		},
		{
			userID:       id2,
			tokenSecret:  secret,
			decodeSecret: secret,
			expiresIn:    duration,
			expected:     id2,
		},
		{
			userID:       id1,
			tokenSecret:  secret,
			decodeSecret: wrongSecret,
			expiresIn:    duration,
			expected:     uuid.UUID{},
		},
		{
			userID:       id1,
			decodeSecret: secret,
			tokenSecret:  secret,
			expiresIn:    expired,
			expected:     uuid.UUID{},
		},

		// add more cases here
	}
	for _, c := range cases {
		tokenString, err := MakeJWT(c.userID, c.tokenSecret, c.expiresIn)
		if err != nil {
			t.Errorf("error from MakeJWT function %sv", err)
		}
		result, err := ValidateJWT(tokenString, c.decodeSecret)
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		//fmt.Printf("result %v, expected %v\n", result, c.expected)

		if result != c.expected {
			t.Errorf("JWT check gave inaccurate answer, actual %v, expected %v", result, c.expected)
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

	}
}
