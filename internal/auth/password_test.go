package auth

import (
	"testing"
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
		input    string
		expected bool
		fakeHash string
	}{
		{
			input:    "george",
			expected: true,
		},
		{
			input:    "george",
			expected: false,
			fakeHash: "lkjldjgekdfgj",
		},

		// add more cases here
	}
	for _, c := range cases {
		hashed, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("error from hash function %sv", err)
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
