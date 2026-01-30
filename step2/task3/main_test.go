package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func removeTestFile(fileName string) error {
	info, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return os.Remove(fileName)
}

func TestLineByNum(t *testing.T) {
	var tests = []struct {
		name           string
		infileName     string
		outFileName    string
		offset         int
		errWasExpected bool
		expected       []byte
		createFile     func(fileName string) error
		readOutFile    func(fileName string) ([]byte, error)
	}{
		{
			name:           "Empty file name",
			infileName:     "",
			outFileName:    "testout.txt",
			offset:         10,
			errWasExpected: true,
			expected:       []byte(nil),
			createFile:     func(fileName string) error { return nil },
			readOutFile:    func(fileName string) ([]byte, error) { return []byte(nil), nil },
		},
		{
			name:           "Empty out file name",
			infileName:     "test.txt",
			outFileName:    "",
			offset:         10,
			errWasExpected: true,
			expected:       []byte(nil),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte(nil), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) { return []byte(nil), nil },
		},
		{
			name:           "Empty file",
			infileName:     "test1.txt",
			outFileName:    "testout1.txt",
			offset:         10,
			errWasExpected: false,
			expected:       []byte(nil),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte(nil), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) {
				return os.ReadFile(fileName)
			},
		},
		{
			name:           "Normal case start with 0",
			infileName:     "test2.txt",
			outFileName:    "testout2.txt",
			offset:         0,
			errWasExpected: false,
			expected:       []byte("abcdefghijklmnopqrstuvwhyz"),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("abcdefghijklmnopqrstuvwhyz"), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) {
				return os.ReadFile(fileName)
			},
		},
		{
			name:           "Normal case start with 13",
			infileName:     "test3.txt",
			outFileName:    "testout3.txt",
			offset:         13,
			errWasExpected: false,
			expected:       []byte("nopqrstuvwhyz"),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("abcdefghijklmnopqrstuvwhyz"), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) {
				return os.ReadFile(fileName)
			},
		},
		{
			name:           "Start less than 0",
			infileName:     "test4.txt",
			outFileName:    "testout4.txt",
			offset:         -1,
			errWasExpected: true,
			expected:       []byte(nil),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("abcdefghijklmnopqrstuvwhyz"), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) {
				return os.ReadFile(fileName)
			},
		},
		{
			name:           "Start greater than file",
			infileName:     "test5.txt",
			outFileName:    "testout5.txt",
			offset:         26,
			errWasExpected: false,
			expected:       []byte(nil),
			createFile: func(fileName string) error {
				return os.WriteFile(fileName, []byte("abcdefghijklmnopqrstuvwhyz"), 0666)
			},
			readOutFile: func(fileName string) ([]byte, error) {
				return os.ReadFile(fileName)
			},
		},
		// Read file error
		// Write file error

	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testDir := os.TempDir()
			testInFile := filepath.Join(testDir, test.infileName)
			testOutFile := filepath.Join(testDir, test.outFileName)
			if err := test.createFile(testInFile); err != nil {
				t.Errorf("cannot create test file: %s\n", err)
			}

			err := CopyFilePart(testInFile, testOutFile, test.offset)
			if (err != nil) != test.errWasExpected {
				t.Errorf(
					"CopyFilePart(\"%s\", \"%s\", %d) got error %v, error was expected %v\n",
					test.infileName,
					test.outFileName,
					test.offset,
					err,
					test.errWasExpected,
				)
			}

			got, err := test.readOutFile(testOutFile)
			if err != nil {
				t.Fatalf("cannot read test out file: %s", err)
			}
			if !bytes.Equal(got, test.expected) {
				t.Errorf(
					"CopyFilePart(\"%s\", \"%s\", %d) got %v, expected %v\n",
					test.infileName,
					test.outFileName,
					test.offset,
					got,
					test.expected,
				)
			}

			if err := removeTestFile(testInFile); err != nil {
				t.Errorf("cannot delete test file: %s\n", err)
			}
			if err := removeTestFile(testOutFile); err != nil {
				t.Errorf("cannot delete test out file: %s\n", err)
			}
		})
	}
}
