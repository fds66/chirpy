package auth

import (
	//"fmt"
	"net/http"
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

// Header map[string][]string
func TestGetBearerToken(t *testing.T) {
	header1 := http.Header{}
	header1.Add("Authorization", "Bearer TOKEN_STRING")
	header2 := http.Header{}
	header2.Add("Authorization", "TOKEN_STRING")
	header3 := http.Header{}
	header3.Add("Content-Type", "application/json")
	header4 := http.Header{}
	header4.Add("Authorization", "")

	cases := []struct {
		input    http.Header
		expected string
		err      bool
	}{
		{
			input:    header1,
			expected: "TOKEN_STRING",
			err:      false,
		},
		{
			input:    header2,
			expected: "",
			err:      true,
		},
		{
			input:    header3,
			expected: "",
			err:      true,
		},
		{
			input:    header4,
			expected: "",
			err:      true,
		},

		// add more cases here
	}
	for _, c := range cases {
		output, outputerr := GetBearerToken(c.input)
		var errorNil bool
		if outputerr != nil {
			errorNil = false
		} else {
			errorNil = true
		}

		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if output != c.expected {
			t.Errorf("Output not as expected, expected %s,%v,  received  %s, %v, %v", c.expected, c.err, output, outputerr, errorNil)
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

	}
}
