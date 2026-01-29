package main

import "testing"

func TestWriteByUpperWriter(t *testing.T) {
	var tests = []struct {
		name           string
		p              []byte
		expectedN      int
		expectedValue  string
		errWasExpected bool
	}{
		{
			name:           "Empty input",
			p:              []byte{},
			expectedN:      0,
			expectedValue:  "",
			errWasExpected: false,
		},
		{
			name:           "Normal input",
			p:              []byte("abcdefghijklmnopqrstuvwxyz"),
			expectedN:      len("abcdefghijklmnopqrstuvwxyz"),
			expectedValue:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			errWasExpected: false,
		},
		{
			name:           "Normal input",
			p:              []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			expectedN:      len("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			expectedValue:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			errWasExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := UpperWriter{}

			gotN, err := writer.Write(test.p)
			if (err != nil) != test.errWasExpected {
				t.Errorf("UpperWriter.Write(%v) got error \"%v\", error was expected \"%v\"\n", test.p, err, test.errWasExpected)
			}
			if gotN != test.expectedN {
				t.Errorf("UpperWriter.Write(%v) got %d, expected %d\n", test.p, gotN, test.expectedN)
			}
			if writer.UpperString != test.expectedValue {
				t.Errorf("UpperWriter.Write(%v) write \"%s\", expected \"%s\"\n", test.p, writer.UpperString, test.expectedValue)
			}
		})
	}
}
