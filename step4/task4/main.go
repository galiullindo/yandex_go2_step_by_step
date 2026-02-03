package main

import "sync"

var (
	Buf   []int
	mutex sync.Mutex
)

func Write(number int) {
	mutex.Lock()
	defer mutex.Unlock()
	Buf = append(Buf, number)
}

func Consume() int {
	mutex.Lock()
	defer mutex.Unlock()

	if len(Buf) > 0 {
		number := Buf[0]
		Buf = Buf[1:]
		return number
	}
	return 0
}
