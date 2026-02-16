package main

import (
	"io"
	"strings"
	"testing"

	"github.com/galiullindo/go-2-step-by-step/step1/testutils"
)

var tests = []struct {
	name      string
	reader    io.Reader
	seq       []byte
	expected  bool
	wantError bool
}{
	{
		name:      "Empty sequence",
		reader:    strings.NewReader(""),
		seq:       []byte(""),
		expected:  false,
		wantError: true,
	},
	{
		name:      "Empty reader",
		reader:    strings.NewReader(""),
		seq:       []byte("abc"),
		expected:  false,
		wantError: false,
	},
	{
		name:      "Reader with sequence",
		reader:    strings.NewReader("ababacabcacbcbc"),
		seq:       []byte("abc"),
		expected:  true,
		wantError: false,
	},
	{
		name:      "Reader without sequence",
		reader:    strings.NewReader("ababacabcacbcbc"),
		seq:       []byte("abcd"),
		expected:  false,
		wantError: false,
	},
	{
		name:      "Read error",
		reader:    testutils.NewCustomReader(),
		seq:       []byte("abc"),
		expected:  false,
		wantError: true,
	},
}

func TestContains(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Contains(test.reader, test.seq)
			if (err != nil) != test.wantError {
				t.Errorf("Contains(reader, seq), error %v, want error %v\n", err, test.wantError)
			}
			if got != test.expected {
				t.Errorf("Contains(reader, seq), got %v, want %v\n", got, test.expected)
			}
		})
	}
}
