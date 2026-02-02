package main

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	var tests = []struct {
		name     string
		counter  Count
		times    int
		expected int
	}{
		{
			name:     "Count to 100",
			counter:  &Counter{},
			times:    100,
			expected: 100,
		},
		{
			name:     "Count to 10000",
			counter:  &Counter{},
			times:    10000,
			expected: 10000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var wg sync.WaitGroup
			for range test.times {
				wg.Add(1)
				go func() {
					defer wg.Done()
					test.counter.Increment()
				}()
			}
			wg.Wait()

			got := test.counter.GetValue()
			if got != test.expected {
				t.Errorf("Got %d, expected %d\n", got, test.expected)
			}
		})
	}
}
