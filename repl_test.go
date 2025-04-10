package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	// ...

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello     world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "BoB",
			expected: []string{"bob"},
		},
		{
			input:    "This is a    new day   is it",
			expected: []string{"this", "is", "a", "new", "day", "is", "it"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length Expected: %v, Actual: %v", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: %v, Actual: %v", expectedWord, word)
			}
		}
	}

}
