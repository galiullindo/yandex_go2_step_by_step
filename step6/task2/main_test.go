package main

import (
	"testing"

	"github.com/galiullindo/yandex_go2_step_by_step/step6/testutils"
)

func TestAverage(t *testing.T) {
	var tests = []struct {
		name            string
		studentNames    []string
		expected        int
		isExpectedError bool
		expectedError   error
	}{
		{
			name:            "Case zero studnet",
			studentNames:    []string{},
			expected:        0,
			isExpectedError: true,
			expectedError:   ErrNamesIsEmpty,
		},
		{
			name:         "Case one studnet",
			studentNames: []string{"Jack50"},
			expected:     50,
		},
		{
			name:         "Case three students",
			studentNames: []string{"Jack50", "John40", "Den10"},
			expected:     33,
		},
		{
			name:            "Case bad request",
			studentNames:    []string{"Jack50", ""},
			expected:        0,
			isExpectedError: true,
			expectedError:   BadRequestError,
		},
		{
			name:            "Case student not found",
			studentNames:    []string{"Jack50", "Barbara25"},
			expected:        0,
			isExpectedError: true,
			expectedError:   StudentNotFoundError,
		},
		{
			name:            "Case internal server error",
			studentNames:    []string{"Jack50", "Barbara25&ise=true"},
			expected:        0,
			isExpectedError: true,
			expectedError:   InternalServerError,
		},
	}

	_, start, stop := testutils.NewServer(":8082")

	go start()
	defer stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Average(test.studentNames)

			if (err != nil) != test.isExpectedError {
				t.Errorf("unexpected error: got %v, error is expected %v\n", err, test.isExpectedError)
			}
			if err != test.expectedError && test.expectedError != nil {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedError)
			}

			if got != test.expected {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.studentNames, got, test.expected)
			}
		})
	}
}
