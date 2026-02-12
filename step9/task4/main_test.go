package main

import "testing"

func TestAreAnagrams(t *testing.T) {
	var tests = []struct {
		name     string
		word     string
		wordToo  string
		expected bool
	}{
		{
			name:     "Case only upper letters",
			word:     "ABCD",
			wordToo:  "DCBA",
			expected: true,
		},
		{
			name:     "Case only lower letters",
			word:     "abcd",
			wordToo:  "dcba",
			expected: true,
		},
		{
			name:     "Case upper and lower letters",
			word:     "AbCd",
			wordToo:  "DcBa",
			expected: true,
		},
		{
			name:     "Case words of different lengths",
			word:     "ABCD",
			wordToo:  "ABCDe",
			expected: false,
		},
		{
			name:     "Case words of the same length",
			word:     "ABCDf",
			wordToo:  "ABCDe",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AreAnagrams(test.word, test.wordToo)
			if got != test.expected {
				t.Errorf("unexpected value for %v, %v: got %v, expected %v\n", test.word, test.wordToo, got, test.expected)
			}
		})
	}
}
