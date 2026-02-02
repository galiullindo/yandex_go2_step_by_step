package main

import "testing"

func TestSend(t *testing.T) {
	channel := make(chan int)
	number := 1
	expected := 1

	go Send(channel, number)
	if got := <-channel; got != expected {
		t.Errorf("go Send(%v, %d) got %d, expected %d\n", channel, number, got, expected)
	}
}
