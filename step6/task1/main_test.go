package main

import (
	"testing"

	"github.com/galiullindo/go-2-step-by-step/step6/testutils"
)

func TestFetchMark(t *testing.T) {
	var tests = []struct {
		name          string
		studentName   string
		expected      Student
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:        "Case normal",
			studentName: "Sara60",
			expected:    Student{Name: "Sara60", Mark: 60},
		},
		{
			name:          "Case error in new request",
			studentName:   "\nAbcd",
			expected:      Student{},
			isExpectedErr: true,
		},
		{
			name:          "Case error in client do",
			studentName:   "Barbara25&abort=true",
			expected:      Student{Name: "Barbara25&abort=true"},
			isExpectedErr: true,
		},
		{
			name:          "Case error in read body",
			studentName:   "Barbara25&read=true",
			expected:      Student{Name: "Barbara25&read=true"},
			isExpectedErr: true,
		},
		{
			name:          "Case error in convert string to integer",
			studentName:   "Barbara25&conv=true",
			expected:      Student{Name: "Barbara25&conv=true"},
			isExpectedErr: true,
		},
		{
			name:          "Case status bad request",
			studentName:   "",
			expected:      Student{},
			isExpectedErr: true,
			expectedErr:   BadRequestError,
		},
		{
			name:          "Case status not found",
			studentName:   "Abcd",
			expected:      Student{Name: "Abcd"},
			isExpectedErr: true,
			expectedErr:   NotFoundError,
		},
		{
			name:          "Case status internal server error",
			studentName:   "Barbara25&ise=true",
			expected:      Student{Name: "Barbara25&ise=true"},
			isExpectedErr: true,
			expectedErr:   InternalServerError,
		},
	}

	_, start, stop := testutils.NewServer(":8082")
	go start()
	defer stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := fetchMark(test.studentName)

			if (err != nil) != test.isExpectedErr {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedErr)
			}
			if test.expectedErr != nil && err != test.expectedErr {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if got != test.expected {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.studentName, got, test.expected)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		name          string
		studentName1  string
		studentName2  string
		expected      string
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:         "Case first student greater than second",
			studentName1: "Jack50",
			studentName2: "John40",
			expected:     ">",
		},
		{
			name:         "Case first student less than second",
			studentName1: "Jack50",
			studentName2: "Sara60",
			expected:     "<",
		},
		{
			name:         "Case first student equal second",
			studentName1: "Jack50",
			studentName2: "Bob50",
			expected:     "=",
		},
		{
			name:         "Case first student is second",
			studentName1: "Jack50",
			studentName2: "Jack50",
			expected:     "=",
		},
		{
			name:          "Case first student with error",
			studentName1:  "",
			studentName2:  "Jack50",
			expected:      "",
			isExpectedErr: true,
		},
		{
			name:          "Case second student with error",
			studentName1:  "Jack50",
			studentName2:  "",
			expected:      "",
			isExpectedErr: true,
		},
	}

	_, start, stop := testutils.NewServer(":8082")
	go start()
	defer stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Compare(test.studentName1, test.studentName2)

			if (err != nil) != test.isExpectedErr {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedErr)
			}
			if test.expectedErr != nil && err != test.expectedErr {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if got != test.expected {
				t.Errorf("unexpected value for %v, %v: got %v, expected %v\n", test.studentName1, test.studentName2, got, test.expected)
			}
		})
	}
}
