package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   AAA BBB         CCCC",
			expected: []string{"aaa", "bbb", "cccc"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check len
		if len(actual) != len(c.expected) {
			t.Errorf("Failed test case:\n input: %v\noutput: %v\nexpected: %v\n", c.input, actual, c.expected)
			t.Fail()
			return
		}

		// Check every word
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Failed test case:\n input: %v\noutput: %v\nexpected: %v\n", c.input, actual, c.expected)
				t.Fail()
				return
			}
		}
	}
}
