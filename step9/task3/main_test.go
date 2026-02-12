package main

import "testing"

func TestReverseString(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"
	expected := "zyxwvutsrqponmlkjihgfedcba"

	got := ReverseString(s)
	if got != expected {
		t.Errorf("unexpected value for %v: got %v, expected %v\n", s, got, expected)
	}
}
