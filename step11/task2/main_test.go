package main

import (
	"reflect"
	"testing"
)

type testCase interface {
	run(t *testing.T)
}

type genericTestCase[T any] struct {
	name      string
	s         []T
	predicate func(T) bool
	expected  []T
}

func (tc *genericTestCase[T]) run(t *testing.T) {
	t.Helper()
	t.Run(tc.name, func(t *testing.T) {
		t.Helper()
		t.Parallel()

		got := Filter(tc.s, tc.predicate)

		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("unexpected value for %v: got %v, expected %v\n", tc.s, got, tc.expected)
		}
	})
}

type custom struct {
	c string
}

func TestSum(t *testing.T) {
	var tests = []testCase{
		&genericTestCase[int]{
			name:      "Case integers: only even ones",
			s:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{2, 4, 6, 8, 0},
		},
		&genericTestCase[float64]{
			name:      "Case floats: only greater than five",
			s:         []float64{1.0, 2.9, 3.8, 4.7, 5.6, 6.5, 7.4, 8.3, 9.2, 0.1},
			predicate: func(f float64) bool { return f > 5 },
			expected:  []float64{5.6, 6.5, 7.4, 8.3, 9.2},
		},
		&genericTestCase[string]{
			name: "Case strings: only those starting with Aa",
			s:    []string{"", "ab", "AB", "ba", "BA"},
			predicate: func(s string) bool {
				if s != "" {
					r := []rune(s)
					return r[0] == 'a' || r[0] == 'A'
				}
				return false
			},
			expected: []string{"ab", "AB"},
		},
		&genericTestCase[custom]{
			name:      "Case castom type: only c equel x",
			s:         []custom{{"a"}, {"b"}, {"c"}, {"x"}, {"x"}, {"d"}},
			predicate: func(c custom) bool { return c.c == "x" },
			expected:  []custom{{"x"}, {"x"}},
		},
	}

	for _, tc := range tests {
		tc.run(t)
	}
}
