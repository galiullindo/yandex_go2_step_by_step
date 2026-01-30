package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadContent(t *testing.T) {
	var tests = []struct {
		name       string
		fileName   string
		expected   string
		createFile func(fileName string) error
	}{
		{
			name:       "Empty file name",
			fileName:   "",
			expected:   "",
			createFile: func(fileName string) error { return nil },
		},
		{
			name:     "Empty file",
			fileName: "test.txt",
			expected: "",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte(nil), 0666)
			},
		},
		{
			name:     "Normal file",
			fileName: "test.txt",
			expected: "abcdefghijklmnopqrstuvwxyz",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("abcdefghijklmnopqrstuvwxyz"), 0666)
			},
		},
		{
			name:       "File not exists",
			fileName:   "test.txt",
			expected:   "",
			createFile: func(fileName string) error { return nil },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testDir := t.TempDir()
			testFile := filepath.Join(testDir, test.fileName)

			if err := test.createFile(testFile); err != nil {
				t.Errorf("cannot create test file\n")
			}

			if got := ReadContent(testFile); got != test.expected {
				t.Errorf("ReadContent(\"%s\") got \"%s\", expected \"%s\"\n", test.fileName, got, test.expected)
			}
		})
	}
}
