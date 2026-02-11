package main

import (
	"slices"
	"testing"
	"time"
)

func TestIsPrime(t *testing.T) {
	var tests = []struct {
		name     string
		numbers  []int
		expected bool
	}{
		{
			name:     "Case prime numbers",
			numbers:  []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47},
			expected: true,
		},
		{
			name:     "Case not prime numbers",
			numbers:  []int{1, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20, 21, 22, 24},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, number := range test.numbers {
				got := isPrime(number)
				if got != test.expected {
					t.Errorf("unexpected value for %v: got %v, expected %v\n", number, got, test.expected)
				}
			}
		})
	}
}

func TestGeneratePrimeNumbers(t *testing.T) {
	// 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47
	var tests = []struct {
		name               string
		n                  int
		expected           []int
		isNotCheckExpected bool
		expectedExecTime   time.Duration
	}{
		{
			name:             "Case n=1",
			n:                1,
			expected:         []int{},
			expectedExecTime: 100 * time.Millisecond,
		},
		{
			name:             "Case n=10",
			n:                10,
			expected:         []int{2, 3, 5, 7},
			expectedExecTime: 100 * time.Millisecond,
		},
		{
			name:             "Case n=11",
			n:                11,
			expected:         []int{2, 3, 5, 7},
			expectedExecTime: 100 * time.Millisecond,
		},
		{
			name:               "Case n=1000000",
			n:                  1000000,
			expected:           []int{},
			isNotCheckExpected: true,
			expectedExecTime:   100 * time.Millisecond,
		},
	}

	for _, test := range tests[3:] {
		t.Run(test.name, func(t *testing.T) {
			stopCh := make(chan struct{})
			primeNumbersCh := make(chan int)

			start := time.Now()

			go GeneratePrimeNumbers(stopCh, primeNumbersCh, test.n)

			got := make([]int, 0)
			for primeNumber := range primeNumbersCh {
				if !test.isNotCheckExpected {
					got = append(got, primeNumber)
				}
			}

			duration := time.Since(start)

			if duration > (test.expectedExecTime + 2*time.Millisecond) {
				t.Errorf(
					"unexpected execution time: got %v ms, expected near %v ms\n",
					duration.Milliseconds(),
					(test.expectedExecTime + 2*time.Millisecond).Milliseconds(),
				)
			}

			if !slices.Equal(got, test.expected) {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.n, got, test.expected)
			}
		})
	}
}
