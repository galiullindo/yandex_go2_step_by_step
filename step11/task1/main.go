package main

type Numbers interface {
	int | float64
}

func Sum[T Numbers](s []T) T {
	var sum T
	for _, item := range s {
		sum += item
	}
	return sum
}
