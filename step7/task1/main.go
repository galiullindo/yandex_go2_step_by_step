package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNegative = errors.New("n cannot be negative")
	ErrTimeout  = errors.New("canceled by timeout")
)

func TimeoutFibonacci(n int, timeout time.Duration) (int, error) {
	timer := time.NewTimer(timeout)

	select {
	case <-timer.C:
		fmt.Println("herer")
		return 0, ErrTimeout
	default:
		if n < 0 {
			return 0, ErrNegative
		}
		switch n {
		case 0:
			return 0, nil
		case 1:
			return 1, nil
		}
	}

	a, b := 0, 1
	for i := 1; i < n; i++ {
		select {
		case <-timer.C:
			return 0, ErrTimeout
		default:
			a, b = b, a+b
		}
	}
	return b, nil
}
