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

func TestAverageMark(t *testing.T) {
	var tests = []struct {
		name     string
		students []Student
		expected int
	}{
		{
			name:     "Case zero students",
			students: []Student{},
			expected: 0,
		},
		{
			name:     "Case one student",
			students: []Student{{Name: "Aaaa", Mark: 25}},
			expected: 25,
		},
		{
			name:     "Case three students",
			students: []Student{{Name: "Aaaa", Mark: 25}, {Name: "Aaaa", Mark: 15}, {Name: "Aaaa", Mark: 33}}, // 24.333...
			expected: 24,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := averageMark(test.students)
			if got != test.expected {
				t.Errorf("unexpected value for %v: got %d, expected %d\n", test.students, got, test.expected)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	var tests = []struct {
		name          string
		studentNames  []string
		expected      int
		isExpectedErr bool
		expectedErr   error
	}{
		{
			name:          "Case zero students",
			studentNames:  []string{},
			expected:      0,
			isExpectedErr: true,
			expectedErr:   NamesIsEmptyError,
		},
		{
			name:         "Case normal students",
			studentNames: []string{"Bob50", "Jack50", "John40", "Den10"}, // avg: 37.5 -> 37
			expected:     37,
		},
		{
			name:          "Case student with not found error",
			studentNames:  []string{"Bob50", "Jack50", "Abcd", "John40", "Den10"},
			expected:      0,
			isExpectedErr: true,
			expectedErr:   NotFoundError,
		},
		{
			name:          "Case student with bad request error",
			studentNames:  []string{"Bob50", "Jack50", "", "John40", "Den10"},
			expected:      0,
			isExpectedErr: true,
			expectedErr:   BadRequestError,
		},
		{
			name:          "Case student with internal server error",
			studentNames:  []string{"Bob50", "Jack50", "Barbara25&ise=true", "John40", "Den10"},
			expected:      0,
			isExpectedErr: true,
			expectedErr:   InternalServerError,
		},
	}

	_, start, stop := testutils.NewServer(":8082")
	go start()
	defer stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Average(test.studentNames)

			if (err != nil) != test.isExpectedErr {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedErr)
			}
			if test.expectedErr != nil && err != test.expectedErr {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if got != test.expected {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.studentNames, got, test.expected)
			}
		})
	}
}
