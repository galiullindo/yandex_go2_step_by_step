package main

import "testing"

func TestSend(t *testing.T) {
	channel1 := make(chan int)
	channel2 := make(chan int)

	Send(channel1, channel2)

	for i := 0; i < 3; i++ {
		expected := i

		if got1 := <-channel1; got1 != expected {
			t.Errorf("Send() got from channel 1 %d, expected %d\n", got1, expected)
		}
		if got2 := <-channel2; got2 != expected {
			t.Errorf("Send() got from channel 2 %d, expected %d\n", got2, expected)
		}
	}
}
