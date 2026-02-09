package main

import (
	"testing"

	"github.com/galiullindo/yandex_go2_step_by_step/step6/testutils"
)

func TestCompare(t *testing.T) {
	var tests = []struct {
		name            string
		studentName1    string
		studentName2    string
		expected        string
		isExpectedError bool
		expectedError   error
	}{
		{
			name:         "Case firs greater than second",
			studentName1: "Jack50",
			studentName2: "John40",
			expected:     ">",
		},
		{
			name:         "Case firs equals second",
			studentName1: "Jack50",
			studentName2: "Bob50",
			expected:     "=",
		},
		{
			name:         "Case firs less than second",
			studentName1: "Jack50",
			studentName2: "Sara60",
			expected:     "<",
		},
		{
			name:         "Case identical names",
			studentName1: "Jack50",
			studentName2: "Jack50",
			expected:     "=",
		},
		{
			name:            "Case student not found",
			studentName1:    "Jack50",
			studentName2:    "Abcd",
			expected:        "",
			isExpectedError: true,
			expectedError:   StudentNotFoundError,
		},
		{
			name:            "Case bad request",
			studentName1:    "Jack50",
			studentName2:    "",
			expected:        "",
			isExpectedError: true,
			expectedError:   BadRequestError,
		},
		{
			name:            "Case internal server error",
			studentName1:    "Jack50",
			studentName2:    "Barbara25&ise=true",
			expected:        "",
			isExpectedError: true,
			expectedError:   InternalServerError,
		},
	}

	_, start, stop := testutils.NewServer(":8082")

	go start()
	defer stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Compare(test.studentName1, test.studentName2)

			if (err != nil) != test.isExpectedError {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedError)
			}
			if err != test.expectedError && test.expectedError != nil {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedError)
			}

			if got != test.expected {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}
