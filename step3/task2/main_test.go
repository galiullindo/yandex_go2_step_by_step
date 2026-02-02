package main

import "testing"

func TestReceive(t *testing.T) {
	channel := make(chan int)
	number := 1
	expected := 1

	go func() { channel <- number }()
	if got := Receive(channel); got != expected {
		t.Errorf("Receive(%v) got %d, expected %d\n", channel, got, expected)
	}
}
