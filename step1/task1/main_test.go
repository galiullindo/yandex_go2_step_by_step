package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type customWriter struct {
}

func NewCustomWriter() *customWriter {
	return &customWriter{}
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

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
			writer:         NewCustomWriter(),
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
			case *customWriter:
				if test.s != "" || test.expected != "" {
					t.Fatal("Error in test: custom writer only for error return check. For checking write use bytes.Buffer.")
				}
			default:
				t.Fatal("Error in test: unexpectable writer.")
			}
		})
	}
}
