package main

import (
	"testing"
)

func TestCleanBody(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input:    "SHARBERT sharbert! sharbert",
			expected: "**** sharbert! ****",
		},
		{
			input:    "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			input:    "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			input:    "I really need a KERFUFFLE to go to bed sooner, Fornax!",
			expected: "I really need a **** to go to bed sooner, Fornax!",
		},
	}

	for _, c := range cases {
		actual := replaceProfaneWords(c.input)
		if actual != c.expected {
			t.Errorf("replaceProfaneWords(%q)==%q, expected %q", c.input, actual, c.expected)
		}
	}

}
