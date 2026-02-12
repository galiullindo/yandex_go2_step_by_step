package main

import (
	"math"
	"testing"
)

func assertEqual[T Numbers](t *testing.T, s []T, got T, expected T) {
	t.Helper()

	if gotF, ok := any(got).(float64); ok {
		expectedF := any(expected).(float64)
		if math.Abs(gotF-expectedF) > 1e-10 {
			t.Errorf("unexpected value for %v: got %v, expected %v\n", s, got, expected)
		}
		return
	}

	if got != expected {
		t.Errorf("unexpected value for %v: got %v, expected %v\n", s, got, expected)
	}
}

type testCase interface {
	run(t *testing.T)
}

type genericTestCase[T Numbers] struct {
	name     string
	s        []T
	expected T
}

func (tc *genericTestCase[T]) run(t *testing.T) {
	t.Helper()
	t.Run(tc.name, func(t *testing.T) {
		t.Helper()
		t.Parallel()

		got := Sum(tc.s)
		assertEqual(t, tc.s, got, tc.expected)
	})
}

func TestSum(t *testing.T) {
	var tests = []testCase{
		// --- integers
		&genericTestCase[int]{
			name:     "Case integers: only zero",
			s:        []int{0, 0, 0, 0, 0},
			expected: 0,
		},
		&genericTestCase[int]{
			name:     "Case integers: only positive",
			s:        []int{1, 2, 3, 4, 5},
			expected: 15,
		},
		&genericTestCase[int]{
			name:     "Case integers: only negative",
			s:        []int{-1, -2, -3, -4, -5},
			expected: -15,
		},
		&genericTestCase[int]{
			name:     "Case integers: negative mixed",
			s:        []int{-1, 2, -3, 4, -5},
			expected: -3,
		},
		&genericTestCase[int]{
			name:     "Case integers: positive mixed",
			s:        []int{1, -2, 3, -4, 5},
			expected: 3,
		},
		// --- floats
		&genericTestCase[float64]{
			name:     "Case float: only zero",
			s:        []float64{0, 0, 0, 0, 0},
			expected: 0,
		},
		&genericTestCase[float64]{
			name:     "Case float: only positive",
			s:        []float64{1.1, 2.2, 3.3, 4.4, 5.5},
			expected: 16.5,
		},
		&genericTestCase[float64]{
			name:     "Case float: only negative",
			s:        []float64{-1.1, -2.2, -3.3, -4.4, -5.5},
			expected: -16.5,
		},
		&genericTestCase[float64]{
			name:     "Case float: negative mixed",
			s:        []float64{-1.1, 2.2, -3.3, 4.4, -5.5},
			expected: -3.3,
		},
		&genericTestCase[float64]{
			name:     "Case float: positive mixed",
			s:        []float64{1.1, -2.2, 3.3, -4.4, 5.5},
			expected: 3.3,
		},
	}

	for _, tc := range tests {
		tc.run(t)
	}
}
