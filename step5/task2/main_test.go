package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func staringTestAPI(t *testing.T) {
	http.HandleFunc(
		"/normal",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("data"))
		},
	)
	http.HandleFunc(
		"/timeout",
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(15 * time.Millisecond)
			w.Write([]byte("data"))
		},
	)
	http.HandleFunc(
		"/error/response",
		func(w http.ResponseWriter, r *http.Request) {
			panic(http.ErrAbortHandler)
		},
	)
	http.HandleFunc(
		"/error/read",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1000")
			w.WriteHeader(http.StatusOK)

			hijacker, _ := w.(http.Hijacker)
			conn, bufrw, _ := hijacker.Hijack()
			_ = bufrw.Flush()
			conn.Close()
		},
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		t.Errorf("cannot staring the test api: %s\n", err)
	}
}

func TestFetchAPI(t *testing.T) {
	var tests = []struct {
		name           string
		ctx            context.Context
		url            string
		timeout        time.Duration
		expected       *APIResponse
		errWasExpected bool
		expectedErr    error
	}{
		{
			name:     "Case normal",
			ctx:      context.Background(),
			url:      "http://localhost:8080/normal",
			timeout:  10 * time.Millisecond,
			expected: &APIResponse{Data: "data", StatusCode: http.StatusOK},
		},
		{
			name:           "Case nil context",
			ctx:            nil,
			url:            "http://localhost:8080/normal",
			timeout:        10 * time.Millisecond,
			expected:       nil,
			errWasExpected: true,
			expectedErr:    ErrNilContext,
		},
		{
			name:     "Case timeout greater than fetch time",
			ctx:      context.Background(),
			url:      "http://localhost:8080/timeout",
			timeout:  20 * time.Millisecond,
			expected: &APIResponse{Data: "data", StatusCode: http.StatusOK},
		},
		{
			name:           "Case timeout less than fetch time",
			ctx:            context.Background(),
			url:            "http://localhost:8080/timeout",
			timeout:        10 * time.Millisecond,
			expected:       nil,
			errWasExpected: true,
			expectedErr:    context.DeadlineExceeded,
		},
		{
			name:           "Case bad url string",
			ctx:            context.Background(),
			url:            "http://localhost:PORT/normal",
			timeout:        10 * time.Millisecond,
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Case response error",
			ctx:            context.Background(),
			url:            "http://localhost:8080/error/response",
			timeout:        10 * time.Millisecond,
			expected:       nil,
			errWasExpected: true,
		},
		{
			name:           "Case read error",
			ctx:            context.Background(),
			url:            "http://localhost:8080/error/read",
			timeout:        10 * time.Millisecond,
			expected:       nil,
			errWasExpected: true,
		},
	}

	go staringTestAPI(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := fetchAPI(test.ctx, test.url, test.timeout)

			if (err != nil) != test.errWasExpected {
				t.Errorf("unexpected error, got %v, err was expected %v\n", err, test.errWasExpected)
			}
			if test.errWasExpected && test.expectedErr != nil {
				if !errors.Is(err, test.expectedErr) {
					t.Errorf("unexpected error, got %v, expected %v\n", err, test.expectedErr)
				}
			}

			isGotNil := got == nil
			isExpectedNil := test.expected == nil
			if !isGotNil && !isExpectedNil {
				if *got != *test.expected {
					t.Errorf("unexpected value, got %v, expected %v", got, test.expected)
				}
			} else {
				if got != test.expected {
					t.Errorf("unexpected value, got %v, expected %v", got, test.expected)
				}
			}
		})
	}
}
