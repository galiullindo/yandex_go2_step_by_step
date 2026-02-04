package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

type CustomReader struct {
	f func(p []byte) (n int, err error)
}

// Attention: if delay greater than 0 copy to dst []byte("abcd") and return 4, io.EOF
func NewCustomReader(f func(p []byte) (n int, err error)) *CustomReader {
	return &CustomReader{f: f}
}

func (r *CustomReader) Read(b []byte) (n int, err error) {
	n, err = r.f(b)
	return n, err
}

func TestContains(t *testing.T) {
	var tests = []struct {
		name           string
		timeout        time.Duration
		reader         io.Reader
		sequence       []byte
		expected       bool
		errWasExpected bool
	}{
		{
			name:           "Case normal",
			timeout:        0,
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			sequence:       []byte("opqr"),
			expected:       true,
			errWasExpected: false,
		},
		{
			name:           "Case value true and channel was closed",
			timeout:        0,
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			sequence:       []byte("wxyz"),
			expected:       true,
			errWasExpected: false,
		},
		{
			name:           "Case value false and channel was closed",
			timeout:        0,
			reader:         bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			sequence:       []byte("wxyza"),
			expected:       false,
			errWasExpected: false,
		},
		{
			name:           "Case empty sequence",
			timeout:        0,
			reader:         NewCustomReader(func(p []byte) (n int, err error) { return 0, io.EOF }),
			sequence:       []byte(nil),
			expected:       false,
			errWasExpected: true,
		},
		{
			name:           "Case read error",
			timeout:        0,
			reader:         NewCustomReader(func(p []byte) (n int, err error) { return 0, errors.New("error") }),
			sequence:       []byte("opqr"),
			expected:       false,
			errWasExpected: true,
		},
		{
			name:    "Case timeout less than delay",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(15 * time.Millisecond)
				n = copy(p, []byte("a"))
				return n, io.EOF
			}),
			sequence:       []byte("a"),
			expected:       false,
			errWasExpected: true,
		},
		{
			name:    "Case value true and timeout greater than delay",
			timeout: 15 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(10 * time.Millisecond)
				n = copy(p, []byte("a"))
				return n, io.EOF
			}),
			sequence:       []byte("a"),
			expected:       true,
			errWasExpected: false,
		},
		{
			name:    "Case value false and timeout greater than delay",
			timeout: 15 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(10 * time.Millisecond)
				n = copy(p, []byte("a"))
				return n, io.EOF
			}),
			sequence:       []byte("b"),
			expected:       false,
			errWasExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.Background(), func() {}
			if test.timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, test.timeout)
			}
			defer cancel()

			got, err := Contains(ctx, test.reader, test.sequence)

			if (err != nil) != test.errWasExpected {
				t.Errorf("unexpected error, got %v, was expected %v\n", err, test.errWasExpected)
			}

			if got != test.expected {
				t.Errorf("unexpected value, got %v, expected %v\n", got, test.expected)
			}
		})
	}
}
