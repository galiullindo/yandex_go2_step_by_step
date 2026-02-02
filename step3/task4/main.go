package main

func Process(numbers []int) chan int {
	newCannel := make(chan int, len(numbers))
	for _, number := range numbers {
		newCannel <- number
	}
	return newCannel
}
