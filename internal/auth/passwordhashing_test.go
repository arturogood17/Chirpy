package auth

import (
	"testing"
)

func Test_HashingPasswords(t *testing.T) {
	cases := []struct {
		input    string
		expected error
	}{
		{
			input:    "password123",
			expected: nil,
		},
		{
			input:    "password",
			expected: nil,
		},
		{
			input:    "prototype13",
			expected: nil,
		},
	}
	for _, test := range cases {
		t.Run(test.input, func(t *testing.T) {
			actual, err := HashPassword(test.input)
			if err != nil {
				t.Errorf("Error when doing the test: %v", err)
			}
			if err = CheckPasswordHash(actual, test.input); err != nil {
				t.Errorf("Test failed. Check password hashing failed: %v", err)
			}
		})
	}
}
