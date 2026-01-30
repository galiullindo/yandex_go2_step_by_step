package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLineByNum(t *testing.T) {
	var tests = []struct {
		name       string
		fileName   string
		lineNumber int
		expected   string
		createFile func(fileName string) error
	}{
		{
			name:       "Empty file name",
			fileName:   "",
			lineNumber: 10,
			expected:   "",
			createFile: func(fileName string) error { return nil },
		},
		{
			name:       "Empty file",
			fileName:   "test.txt",
			lineNumber: 10,
			expected:   "",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte(nil), 0666)
			},
		},
		{
			name:       "Line number is less than zero",
			fileName:   "test2.txt",
			lineNumber: -1,
			expected:   "",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("0\n1ab\n2cd\n3ef\n4gh\n5ij\n6kl\n7mn\n8op\n9qr\n10st\n11uv\n12wx\n13yz\n"), 0666)
			},
		},
		{
			name:       "Line number is zero",
			fileName:   "test3.txt",
			lineNumber: 0,
			expected:   "0",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("0\n1ab\n2cd\n3ef\n4gh\n5ij\n6kl\n7mn\n8op\n9qr\n10st\n11uv\n12wx\n13yz\n"), 0666)
			},
		},
		{
			name:       "Normal line number 1",
			fileName:   "test4.txt",
			lineNumber: 1,
			expected:   "1ab",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("0\n1ab\n2cd\n3ef\n4gh\n5ij\n6kl\n7mn\n8op\n9qr\n10st\n11uv\n12wx\n13yz\n"), 0666)
			},
		},
		{
			name:       "Normal line number 10",
			fileName:   "test5.txt",
			lineNumber: 10,
			expected:   "10st",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("0\n1ab\n2cd\n3ef\n4gh\n5ij\n6kl\n7mn\n8op\n9qr\n10st\n11uv\n12wx\n13yz\n"), 0666)
			},
		},
		{
			name:       "Line number is greater than file",
			fileName:   "test6.txt",
			lineNumber: 14,
			expected:   "",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("0\n1ab\n2cd\n3ef\n4gh\n5ij\n6kl\n7mn\n8op\n9qr\n10st\n11uv\n12wx\n13yz\n"), 0666)
			},
		},
		{
			name:       "Scaner error",
			fileName:   "test7.txt",
			lineNumber: 14,
			expected:   "",
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte(strings.Repeat("a", 64*1024+1)), 0666)
			},
		},
		{
			name:       "File not exists",
			fileName:   "test8.txt",
			lineNumber: 10,
			expected:   "",
			createFile: func(fileName string) error { return nil },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testDir := os.TempDir()
			testFile := filepath.Join(testDir, test.fileName)
			if err := test.createFile(testFile); err != nil {
				t.Errorf("cannot create test file: %s\n", err)
			}

			if got := LineByNum(testFile, test.lineNumber); got != test.expected {
				t.Errorf("LineByNum(\"%s\", %d) got %#v, expected %#v\n", test.fileName, test.lineNumber, got, test.expected)
			}

			info, err := os.Stat(testFile)
			if err == nil {
				if !info.IsDir() {
					if err := os.Remove(testFile); err != nil {
						t.Errorf("cannot delete test file: %s\n", err)
					}
				}

			}
		})
	}
}
