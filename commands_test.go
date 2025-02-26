package main

import "testing"

func TestCommands(t *testing.T) {
	// Setup of test case structs
	/* cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    " this                  is just whitespace                           ",
			expected: []string{"this", "is", "just", "whitespace"},
		},
		{
			input:    "\t\t\t\t\t\t\t\tTEXT\r\n\r\n",
			expected: []string{"text"},
		},
	} */

	// Loop over cases and run the tests
	/* for _, c := range cases {
		actual := cleanInput(c.input)
		// Check length of the actual slice
		// If they don't match, use the t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual (%v) does not equal expected (%v)", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected: %v; Got: %v", expectedWord, word)
			}
		}
	} */

}
