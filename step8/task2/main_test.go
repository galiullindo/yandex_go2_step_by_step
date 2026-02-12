package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

var CustomReadError = errors.New("custom read error")

type CustomReader struct {
	f func(p []byte) (n int, err error)
}

func NewCustomReader(f func(p []byte) (n int, err error)) *CustomReader {
	return &CustomReader{f: f}
}

func (r *CustomReader) Read(b []byte) (n int, err error) {
	n, err = r.f(b)
	return n, err
}

func TestReadAllWithContext(t *testing.T) {
	var tests = []struct {
		name          string
		timeout       time.Duration
		reader        io.Reader
		expected      []byte
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:     "Case normal",
			timeout:  10 * time.Millisecond,
			reader:   bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz")),
			expected: []byte("abcdefghijklmnopqrstuvwxyz"),
		},
		{
			name:     "Case add buffer capacity",
			timeout:  10 * time.Millisecond,
			reader:   bytes.NewReader([]byte(strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100))),
			expected: []byte(strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100)),
		},
		{
			name:    "Case read error",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("abcdefghijklmnopqrstuvwxyz"))
				err = CustomReadError
				return n, err
			}),
			expected:      []byte("abcdefghijklmnopqrstuvwxyz"),
			isExpectedErr: true,
			expectedErr:   CustomReadError,
		},
		{
			name:    "Case endless reading",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = 0
				err = nil
				return n, err
			}),
			expected:      []byte(nil),
			isExpectedErr: true,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:    "Case big delay",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(100 * time.Millisecond)
				n = copy(p, []byte("abcdefghijklmnopqrstuvwxyz"))
				err = io.EOF
				return n, err
			}),
			expected:      []byte(nil),
			isExpectedErr: true,
			expectedErr:   context.DeadlineExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			start := time.Now()
			got, err := ReadAllWithContext(ctx, test.reader)

			duration := time.Since(start)
			if duration > test.timeout+5*time.Millisecond {
				t.Errorf("unexpected execution time: got %v, expected near %v\n", duration, test.timeout)
			}

			if (err != nil) != test.isExpectedErr {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedErr)
			}
			if test.expectedErr != nil && err != test.expectedErr {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if !bytes.Equal(got, test.expected) {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}

func TestMakeChannelForReadeing(t *testing.T) {
	var tests = []struct {
		name        string
		timeout     time.Duration
		filePath    string
		fileContent []byte
		expected    []byte
	}{
		{
			name:        "Case normal",
			timeout:     10 * time.Millisecond,
			filePath:    "test*.json",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    []byte("abcdefghijklmnopqrstuvwxyz"),
		},
		{
			name:        "Case error opening file",
			timeout:     10 * time.Millisecond,
			filePath:    "",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    nil,
		},
		{
			name:        "Case error reading file",
			timeout:     1 * time.Millisecond,
			filePath:    "test*.json",
			fileContent: []byte(strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100)),
			expected:    nil,
		},
		{
			name:        "Case timeout zero",
			timeout:     0,
			filePath:    "test*.json",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFilePath := test.filePath
			if test.filePath != "" {
				tempFile, err := os.CreateTemp("", test.filePath)
				if err != nil {
					t.Fatalf("error creating test file: %v\n", err)
				}
				defer os.Remove(tempFile.Name())

				if _, err := tempFile.Write(test.fileContent); err != nil {
					t.Fatalf("error writeing in test file: %v\n", err)
				}
				if err := tempFile.Close(); err != nil {
					t.Fatalf("error closing test file: %v\n", err)
				}
				tempFilePath = tempFile.Name()
			}

			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			start := time.Now()
			channel := MakeChannelForReading(ctx, tempFilePath)

			timeout := test.timeout + 500*time.Millisecond
			select {
			case <-time.After(timeout):
				t.Errorf("canceled by global timeout %v\n", timeout)
			case got, ok := <-channel:
				if !ok {
					t.Errorf("channel unexpectedly closed\n")
				} else {
					duration := time.Since(start)
					if duration > test.timeout+10*time.Millisecond {
						t.Errorf("unexpected execution time: got %v, expected near %v\n", duration, test.timeout)
					}
					if !bytes.Equal(got, test.expected) {
						t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
					}
				}
			}
		})
	}
}

func TestReadJSON(t *testing.T) {
	var tests = []struct {
		name                  string
		timeout               time.Duration
		filePath              string
		fileContent           []byte
		isNotExpectedTheValue bool
		expected              []byte
	}{
		{
			name:        "Case normal",
			timeout:     10 * time.Millisecond,
			filePath:    "test*.json",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    []byte("abcdefghijklmnopqrstuvwxyz"),
		},
		{
			name:                  "Case error opening file",
			timeout:               10 * time.Millisecond,
			filePath:              "",
			fileContent:           []byte("abcdefghijklmnopqrstuvwxyz"),
			isNotExpectedTheValue: true,
		},
		{
			name:                  "Case error reading file",
			timeout:               1 * time.Millisecond,
			filePath:              "test*.json",
			fileContent:           []byte(strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100)),
			isNotExpectedTheValue: true,
		},
		{
			name:                  "Case timeout zero",
			timeout:               0,
			filePath:              "test*.json",
			fileContent:           []byte("abcdefghijklmnopqrstuvwxyz"),
			isNotExpectedTheValue: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempFilePath := test.filePath
			if test.filePath != "" {
				tempFile, err := os.CreateTemp("", test.filePath)
				if err != nil {
					t.Fatalf("error creating test file: %v\n", err)
				}
				defer os.Remove(tempFile.Name())

				if _, err := tempFile.Write(test.fileContent); err != nil {
					t.Fatalf("error writeing in test file: %v\n", err)
				}
				if err := tempFile.Close(); err != nil {
					t.Fatalf("error closing test file: %v\n", err)
				}
				tempFilePath = tempFile.Name()
			}

			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			channel := make(chan []byte)

			start := time.Now()
			go ReadJSON(ctx, tempFilePath, channel)

			timeout := test.timeout + 500*time.Millisecond
			select {
			case <-time.After(timeout):
				t.Errorf("canceled by global timeout %v\n", timeout)
			case got, ok := <-channel:
				if ok != !test.isNotExpectedTheValue {
					t.Errorf("unexpected value: got %v, %v, the value is not expected %v\n", got, ok, test.isNotExpectedTheValue)
				}
				if ok && !test.isNotExpectedTheValue {
					duration := time.Since(start)
					if duration > test.timeout+10*time.Millisecond {
						t.Errorf("unexpected execution time: got %v, expected near %v\n", duration, test.timeout)
					}
					if !bytes.Equal(got, test.expected) {
						t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
					}
				}
			}
		})
	}
}
