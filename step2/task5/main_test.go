package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestExtractLog(t *testing.T) {
	var tests = []struct {
		name           string
		fileName       string
		start          time.Time
		end            time.Time
		fileContent    []byte
		expected       []string
		errWasExpected bool
	}{
		{
			name:           "Normal file",
			fileName:       "test1.txt",
			start:          time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    []byte("01.01.2026 info\n02.01.2026 info\n03.01.2026 info\n04.01.2026 info\n05.01.2026 info\n"),
			expected:       []string{"02.01.2026 info", "03.01.2026 info", "04.01.2026 info"},
			errWasExpected: false,
		},
		{
			name:           "Normal file with invalid time period error",
			fileName:       "test2.txt",
			start:          time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    []byte("01.01.2026 info\n02.01.2026 info\n03.01.2026 info\n04.01.2026 info\n05.01.2026 info\n"),
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Normal file with bad log extraction error",
			fileName:       "test3.txt",
			start:          time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    []byte("01.01.2026 info\n02.01.2026 info\n03.01.2026 info\n04.01.2026 info\n05.01.2026 info\n"),
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Empty file with bad log extraction error",
			fileName:       "test4.txt",
			start:          time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    nil,
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Empty file name",
			fileName:       "",
			start:          time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    nil,
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Normal file with bad data",
			fileName:       "test5.txt",
			start:          time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			end:            time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC),
			fileContent:    []byte("01.01.2026 info\n02.01.2026 info\n03.0x.2026 info\n04.01.2026 info\n05.01.2026 info\n"),
			expected:       []string{"02.01.2026 info", "04.01.2026 info"},
			errWasExpected: false,
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

			got, err := ExtractLog(testFile, test.start, test.end)
			if (err != nil) != test.errWasExpected {
				t.Errorf(
					"ExtractLog(%v, %v, %v) got error %v, error was expected %v\n",
					test.fileName,
					test.start,
					test.end,
					err,
					test.errWasExpected,
				)
			}
			if !slices.Equal(got, test.expected) {
				t.Errorf(
					"ExtractLog(%v, %v, %v) got %v, expected %v\n",
					test.fileName,
					test.start,
					test.end,
					got,
					test.expected,
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
