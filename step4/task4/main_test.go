package main

import (
	"sync"
	"testing"
)

func TestWriteAndConsume(t *testing.T) {
	var tests = []struct {
		name      string
		values    []int
		expecteds []int
	}{
		{
			name:      "Write and Consume case normal",
			values:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expecteds: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:      "Write and Consume case empty buffer",
			values:    []int(nil),
			expecteds: []int{0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Buf = make([]int, 0)

			for _, value := range test.values {
				Write(value)
			}

			for _, expected := range test.expecteds {
				if got := Consume(); got != expected {
					t.Errorf("Got %d, expected %d\n", got, expected)
				}
			}
		})
	}
}

func TestWriteAndConsumeAsGoroutines(t *testing.T) {
	var tests = []struct {
		name       string
		goroutines int
		expected   int
	}{
		{
			name:       "Parallel processing case 100 goroutines",
			goroutines: 100,
			expected:   100,
		},
		{
			name:       "Parallel processing case 1000 goroutines",
			goroutines: 1000,
			expected:   1000,
		},
		{
			name:       "Parallel processing case 10000 goroutines",
			goroutines: 10000,
			expected:   10000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Buf = make([]int, 0)
			var wg sync.WaitGroup

			for range test.goroutines {
				wg.Go(func() {
					Write(1)
				})
			}
			wg.Wait()

			if got := len(Buf); got != test.expected {
				t.Errorf("unexpected buffer length: got %d, expected %d\n", got, test.expected)
			}

			for range test.goroutines {
				wg.Go(func() {
					Consume()
				})
			}
			wg.Wait()

			if got := len(Buf); got != 0 {
				t.Errorf("unexpected buffer length: got %d, expected 0\n", got)
			}
		})
	}
}
