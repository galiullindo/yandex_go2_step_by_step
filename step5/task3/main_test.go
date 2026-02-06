package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

func NewAPI(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/normal/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "1")
	})
	mux.HandleFunc("/normal/2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", "2")
	})
	mux.HandleFunc("/timeout/time/ms/15", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Millisecond)
		fmt.Fprintf(w, "%s", "timeout1")
	})
	mux.HandleFunc("/timeout/time/ms/5", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		fmt.Fprintf(w, "%s", "timeout2")
	})
	mux.HandleFunc("/errors/response", func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})
	mux.HandleFunc("/errors/read", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		hijacker, _ := w.(http.Hijacker)
		connection, bufferRW, _ := hijacker.Hijack()
		_ = bufferRW.Flush()
		connection.Close()
	})

	return &http.Server{Addr: addr, Handler: mux}
}

func TestFetch(t *testing.T) {
	API := NewAPI(":8080")

	go func() {
		if err := API.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error in the test API: %s\n", err)
		}
	}()

	defer func() {
		if err := API.Close(); err != nil {
			t.Errorf("error stopping test API %s\n", err)
		}
	}()

	var tests = []struct {
		name            string
		url             string
		timeout         time.Duration
		isExpectedError bool
		expected        APIResponse
	}{
		{
			name:    "Case normal 1",
			url:     "http://localhost:8080/normal/1",
			timeout: 10 * time.Millisecond,
			expected: APIResponse{
				URL:        "http://localhost:8080/normal/1",
				Data:       "1",
				StatusCode: http.StatusOK,
			},
		},
		{
			name:    "Case normal 2",
			url:     "http://localhost:8080/normal/2",
			timeout: 10 * time.Millisecond,
			expected: APIResponse{
				URL:        "http://localhost:8080/normal/2",
				Data:       "2",
				StatusCode: http.StatusOK,
			},
		},
		{
			name:            "Case time 15ms with timeout 10ms",
			url:             "http://localhost:8080/timeout/time/ms/15",
			timeout:         10 * time.Millisecond,
			isExpectedError: true,
			expected: APIResponse{
				URL: "http://localhost:8080/timeout/time/ms/15",
				Err: context.DeadlineExceeded,
			},
		},
		{
			name:    "Case time 5ms with timeout 10ms",
			url:     "http://localhost:8080/timeout/time/ms/5",
			timeout: 10 * time.Millisecond,
			expected: APIResponse{
				URL:        "http://localhost:8080/timeout/time/ms/5",
				Data:       "timeout2",
				StatusCode: http.StatusOK,
			},
		},
		{
			name:            "Case error while request",
			url:             "http://localhost:80xx/errors/request",
			timeout:         10 * time.Millisecond,
			isExpectedError: true,
			expected: APIResponse{
				URL: "http://localhost:80xx/errors/request",
			},
		},
		{
			name:            "Case error while response",
			url:             "http://localhost:8080/errors/response",
			timeout:         10 * time.Millisecond,
			isExpectedError: true,
			expected: APIResponse{
				URL: "http://localhost:8080/errors/response",
				Err: io.EOF,
			},
		},
		{
			name:            "Case error while read",
			url:             "http://localhost:8080/errors/read",
			timeout:         10 * time.Millisecond,
			isExpectedError: true,
			expected: APIResponse{
				URL:        "http://localhost:8080/errors/read",
				StatusCode: http.StatusOK,
				Err:        io.ErrUnexpectedEOF,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			response := Fetch(ctx, test.url)

			ok := response.URL == test.expected.URL && response.Data == test.expected.Data && response.StatusCode == test.expected.StatusCode
			if !ok {
				t.Errorf("unexpected response: got %v, expected %v\n", response, test.expected)
			}

			if (response.Err != nil) != test.isExpectedError {
				t.Errorf("unexpected error: got %v, expected %v\n", response.Err, test.expected.Err)
			}
			if test.expected.Err != nil && !errors.Is(response.Err, test.expected.Err) {
				t.Errorf("unexpected error: got %v, expected %v\n", response.Err, test.expected.Err)
			}
		})
	}
}

func TestFetchAPI(t *testing.T) {
	API := NewAPI(":8080")

	go func() {
		if err := API.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error in the test API: %s\n", err)
		}
	}()

	defer func() {
		if err := API.Close(); err != nil {
			t.Errorf("error stopping test API %s\n", err)
		}
	}()

	var expecteds = map[string]APIResponse{
		"http://localhost:8080/normal/1": {
			URL:        "http://localhost:8080/normal/1",
			Data:       "1",
			StatusCode: http.StatusOK,
		},
		"http://localhost:8080/normal/2": {
			URL:        "http://localhost:8080/normal/2",
			Data:       "2",
			StatusCode: http.StatusOK,
		},
		"http://localhost:8080/timeout/time/ms/15": {
			URL: "http://localhost:8080/timeout/time/ms/15",
			Err: context.DeadlineExceeded,
		},
		"http://localhost:8080/timeout/time/ms/5": {
			URL:        "http://localhost:8080/timeout/time/ms/5",
			Data:       "timeout2",
			StatusCode: http.StatusOK,
		},
	}

	var tests = []struct {
		name    string
		urls    []string
		timeout time.Duration
	}{
		{
			name: "Case all responses are normal",
			urls: []string{
				"http://localhost:8080/normal/1",
				"http://localhost:8080/normal/1",
				"http://localhost:8080/normal/2",
				"http://localhost:8080/normal/2",
			},
			timeout: 10 * time.Millisecond,
		},
		{
			name: "Case has responses with timeout error",
			urls: []string{
				"http://localhost:8080/normal/1",
				"http://localhost:8080/normal/1",
				"http://localhost:8080/timeout/time/ms/15",
				"http://localhost:8080/normal/2",
				"http://localhost:8080/normal/2",
				"http://localhost:8080/timeout/time/ms/15",
			},
			timeout: 10 * time.Millisecond,
		},
		{
			name: "Case has responses with timeout less time",
			urls: []string{
				"http://localhost:8080/normal/1",
				"http://localhost:8080/normal/1",
				"http://localhost:8080/timeout/time/ms/15",
				"http://localhost:8080/normal/2",
				"http://localhost:8080/normal/2",
				"http://localhost:8080/timeout/time/ms/15",
				"http://localhost:8080/timeout/time/ms/5",
			},
			timeout: 10 * time.Millisecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			responses := FetchAPI(ctx, test.urls, test.timeout)

			numberOfResponses, numberOfUrls := len(responses), len(test.urls)
			if numberOfResponses != numberOfUrls {
				t.Errorf(
					"number of responses not equals number of urls: got %d, expected %d\n",
					numberOfResponses,
					numberOfUrls,
				)
			}

			for _, response := range responses {
				expected, ok := expecteds[response.URL]
				if !ok {
					t.Fatalf("unexpected url: %s\n", response.URL)
				}

				ok = response.URL == expected.URL && response.Data == expected.Data && response.StatusCode == expected.StatusCode
				if !ok {
					t.Errorf("unexpected response: got %v, expected %v\n", response, expected)
				}
				if !errors.Is(response.Err, expected.Err) {
					t.Errorf("unexpected error: got %v, expected %v\n", response.Err, expected.Err)
				}
			}
		})
	}
}
