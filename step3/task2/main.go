package main

func Receive(channel chan int) int {
	number := <-channel
	return number
}
