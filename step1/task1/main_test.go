package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/galiullindo/go-2-step-by-step/step1/testutils"
)

func TestWriteString(t *testing.T) {
	var tests = []struct {
		name           string
		s              string
		writer         io.Writer
		expected       string
		errWasExpected bool
	}{
		{
			name:           "Empty string",
			s:              "",
			writer:         bytes.NewBuffer([]byte(nil)),
			expected:       "",
			errWasExpected: false,
		},
		{
			name:           "Normal string",
			s:              "abcdefghijklmnopqrstuvwxyz",
			writer:         bytes.NewBuffer([]byte(nil)),
			expected:       "abcdefghijklmnopqrstuvwxyz",
			errWasExpected: false,
		},
		{
			name:           "Write error",
			writer:         testutils.NewCustomWriter(),
			errWasExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WriteString(test.s, test.writer)
			if (err != nil) != test.errWasExpected {
				t.Errorf("WriteString(%s, %v), got error \"%v\", error was expected \"%v\" ", test.s, test.writer, err, test.errWasExpected)
			}

			switch writer := test.writer.(type) {
			case *bytes.Buffer:
				if got := writer.String(); got != test.expected {
					t.Errorf("WriteString(%s, %v), got \"%s\", expected \"%s\" ", test.s, test.writer, got, test.expected)
				}
			case *testutils.CustomWriter:
				if test.s != "" || test.expected != "" {
					t.Fatal("Error in test: custom writer only for error return check. For checking write use bytes.Buffer.")
				}
			default:
				t.Fatal("Error in test: unexpectable writer.")
			}
		})
	}
}
