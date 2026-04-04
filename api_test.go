package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "so much kerfuffle",
			expected: "so much ****",
		},
		{
			input:    "well Sharbert",
			expected: "well ****",
		},
		{
			input:    "got to be fornax",
			expected: "got to be ****",
		},
		{
			input:    "well Sharbert!",
			expected: "well Sharbert!",
		},
		{
			input:    " well Sharbert ",
			expected: " well **** ",
		},

		// add more cases here
	}
	for _, c := range cases {
		actual := cleanString(c.input)

		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if actual != c.expected {
			t.Errorf("words don't match, actual %s, expected %s", actual, c.expected)
		}
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

	}
}
