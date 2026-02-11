package main

import (
	"math"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < int(math.Sqrt(float64(n)))+1; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func GeneratePrimeNumbers(stop chan struct{}, primeNumbers chan int, n int) {
	defer close(primeNumbers)
	timeout := time.AfterFunc(100*time.Millisecond, func() { stop <- struct{}{} })
	defer timeout.Stop()

	number := 0
	for {
		select {
		case <-stop:
			return
		default:
			if number >= n {
				return
			}

			primalityCheck := func() <-chan bool {
				channel := make(chan bool)
				go func() {
					defer close(channel)
					channel <- isPrime(number)
				}()
				return channel
			}()

			select {
			case <-stop:
				return
			case isPrime := <-primalityCheck:
				if isPrime {
					primeNumbers <- number
				}
				number++
			}
		}
	}
}
