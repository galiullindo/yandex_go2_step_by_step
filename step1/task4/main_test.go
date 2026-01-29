package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type customReader struct {
}

func NewCustomReader() *customReader {
	return &customReader{}
}

func (r *customReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

type customWriter struct {
}

func NewCustomWriter() *customWriter {
	return &customWriter{}
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestCopy(t *testing.T) {
	var tests = []struct {
		name           string
		reader         io.Reader
		writer         io.Writer
		n              uint
		expectedValue  string
		errWasExpected bool
	}{
		{
			name:           "Empty reader",
			reader:         strings.NewReader(""),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              10,
			expectedValue:  "",
			errWasExpected: false,
		},
		{
			name:           "Normal input need 0 byte",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              0,
			expectedValue:  "",
			errWasExpected: false,
		},
		{
			name:           "Normal input need 10 bytes",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              10,
			expectedValue:  "abcdefghij",
			errWasExpected: false,
		},
		{
			name:           "Normal input need 20 bytes",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              20,
			expectedValue:  "abcdefghijklmnopqrst",
			errWasExpected: false,
		},
		{
			name:           "Normal input need 30 bytes",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              30,
			expectedValue:  "abcdefghijklmnopqrstuvwxyz",
			errWasExpected: false,
		},
		{
			name:           "Normal input need 100 bytes",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              100,
			expectedValue:  "abcdefghijklmnopqrstuvwxyz",
			errWasExpected: false,
		},
		{
			name:           "Read error",
			reader:         NewCustomReader(),
			writer:         bytes.NewBuffer([]byte(nil)),
			n:              100,
			expectedValue:  "",
			errWasExpected: true,
		},
		{
			name:           "Write error",
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			writer:         NewCustomWriter(),
			n:              100,
			expectedValue:  "",
			errWasExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Copy(test.reader, test.writer, test.n)
			if (err != nil) != test.errWasExpected {
				t.Errorf(
					"Copy(%v, %v, %d) got error \"%v\", error was expected \"%v\"\n",
					test.reader,
					test.writer,
					test.n,
					err,
					test.errWasExpected,
				)
			}

			switch writer := test.writer.(type) {
			case *bytes.Buffer:
				if got := writer.String(); got != test.expectedValue {
					t.Errorf(
						"Copy(%v, %v, %d) write value \"%s\", expected \"%s\"\n",
						test.reader,
						test.writer,
						test.n,
						got,
						test.expectedValue,
					)
				}
			case *customWriter:
				if test.expectedValue != "" {
					t.Fatal("Error in test: custom writer only for error return check. For checking write use bytes.Buffer.")
				}
			default:
				t.Fatal("Error in test: unexpectable writer.")
			}
		})
	}
}
