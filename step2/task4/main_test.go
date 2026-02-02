package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestModifyFile(t *testing.T) {
	var tests = []struct {
		name        string
		fileName    string
		offset      int
		content     string
		fileContent []byte
		expected    []byte
	}{
		{
			name:        "Normal file, offset 0",
			fileName:    "test1.txt",
			offset:      0,
			content:     "vwxyz",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    []byte("vwxyzfghijklmnopqrstuvwxyz"),
		},
		{
			name:        "Normal file, offset 10",
			fileName:    "test2.txt",
			offset:      10,
			content:     "vwxyz",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    []byte("abcdefghijvwxyzpqrstuvwxyz"),
		},
		{
			name:        "Empty file, offset 0",
			fileName:    "test3.txt",
			offset:      0,
			content:     "vwxyz",
			fileContent: []byte(nil),
			expected:    []byte("vwxyz"),
		},
		{
			name:        "Empty file, offset 10",
			fileName:    "test4.txt",
			offset:      10,
			content:     "vwxyz",
			fileContent: []byte(nil),
			expected:    append([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, []byte("vwxyz")...),
		},
		{
			name:        "Empty file name",
			fileName:    "",
			offset:      0,
			content:     "vwxyz",
			fileContent: []byte(nil),
			expected:    []byte(nil),
		},
		{
			name:        "Normal file, offset -1",
			fileName:    "test5.txt",
			offset:      -1,
			content:     "vwxyz",
			fileContent: []byte("abcdefghijklmnopqrstuvwxyz"),
			expected:    []byte("abcdefghijklmnopqrstuvwxyz"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testDir := os.TempDir()
			testFile := filepath.Join(testDir, test.fileName)

			if test.fileName != "" {
				if err := os.WriteFile(testFile, test.fileContent, 0666); err != nil {
					t.Fatalf("cannot create test file %s\n", err)
				}
			}

			ModifyFile(testFile, test.offset, test.content)

			var got, err = []byte(nil), error(nil)
			if test.fileName != "" {
				got, err = os.ReadFile(testFile)
				if err != nil {
					t.Fatalf("cannot read test file %s\n", err)
				}
			}

			if !bytes.Equal(got, test.expected) {
				t.Errorf(
					"ModifyFile(%v, %d, %v) write in file %v, expected %v\n",
					test.fileName,
					test.offset,
					test.content,
					string(got),
					string(test.expected),
				)
			}

			if test.fileName != "" {
				if err := os.Remove(testFile); err != nil {
					t.Errorf("cannot remove test file %s\n", err)
				}
			}
		})
	}
}
