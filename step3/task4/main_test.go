package main

import (
	"slices"
	"testing"
)

func TestProcess(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	channel := Process(numbers)

	got := make([]int, 0, len(numbers))
	for i := 0; i < len(numbers); i++ {
		got = append(got, <-channel)
	}

	if !slices.Equal(got, expected) {
		t.Errorf("Process(%v) got %v, expected %v\n", numbers, got, expected)
	}
}
