package main

import (
	"io"
	"strings"
	"testing"

	"github.com/galiullindo/go-2-step-by-step/step1/testutils"
)

func TestReadString(t *testing.T) {
	var tests = []struct {
		name           string
		reader         io.Reader
		expected       string
		errWasExpected bool
	}{
		{
			name:     "Empty reader",
			reader:   strings.NewReader(""),
			expected: "",
		},
		{
			name:     "Normal reader",
			reader:   strings.NewReader("abcdefg"),
			expected: "abcdefg",
		},
		{
			name:           "Read error",
			reader:         testutils.NewCustomReader(),
			errWasExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ReadString(test.reader)
			if (err != nil) != test.errWasExpected {
				t.Errorf("ReadString(%v) got error \"%v\", error was expected \"%v\"\n", test.reader, err, test.errWasExpected)
			}
			if got != test.expected {
				t.Errorf("ReadString(%v) got \"%s\", expected \"%v\"\n", test.reader, got, test.expected)
			}
		})
	}
}
