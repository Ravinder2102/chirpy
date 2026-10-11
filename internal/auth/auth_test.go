package auth

import (
	"testing"
)

func TestHashing(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "abc",
			expected: true,
		},
		{
			input:    "abc@123ABC",
			expected: true,
		},
	}

	for _, c := range cases {
		hash, _ := HashPassword(c.input)
		match, _ := CheckPasswordHash(c.input, hash)
		if !match {
			t.Errorf("password and hash does not match")
		}
	}

}
