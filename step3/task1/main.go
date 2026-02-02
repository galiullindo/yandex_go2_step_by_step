package main

func Send(channel chan int, number int) {
	channel <- number
}
