package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
)

type Params struct {
	Name string
	ISE  bool
}

func ParseParams(r *http.Request) (Params, error) {
	query := r.URL.Query()
	params := Params{}

	params.Name = query.Get("name")
	if params.Name == "" {
		return params, fmt.Errorf("missing name")
	}

	iseStr := query.Get("ise")
	if iseStr == "" {
		params.ISE = false
	} else {
		ise, err := strconv.ParseBool(iseStr)
		if err != nil {
			params.ISE = false
		} else {
			params.ISE = ise
		}
	}

	return params, nil
}

func NewTestServer(addr string) *http.Server {
	studentMap := map[string]int{
		"Jack50": 50,
		"John40": 40,
		"Bob10":  10,
		"Sara60": 60,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mark", func(w http.ResponseWriter, r *http.Request) {
		params, err := ParseParams(r)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		if params.ISE {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		mark, found := studentMap[params.Name]
		if !found {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "%d", mark)
	})

	return &http.Server{Addr: addr, Handler: mux}
}

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
			studentNames: []string{"Jack50", "John40", "Bob10"},
			expected:     33,
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

	testServer := NewTestServer(":8082")

	go func() {
		if err := testServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("error in the test server: %s\n", err)
		}
	}()

	defer func() {
		if err := testServer.Close(); err != nil {
			log.Printf("error stopping the test server: %s\n", err)
		}
	}()

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
