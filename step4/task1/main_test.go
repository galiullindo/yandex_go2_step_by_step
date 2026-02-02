package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	goroutines := 10

	concurrentMap := NewSafeMap()

	var wg sync.WaitGroup

	// Проверка на отсутсвие race condition
	// race condition - одновременное выполнение операции чтения и записи несколькими потоками.
	for i := range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			concurrentMap.Set(strconv.Itoa(i), i)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			concurrentMap.Get(strconv.Itoa(i))
		}()
	}
	wg.Wait()

	// Проверка на корректность возвращаемых значений.
	for i := range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			concurrentMap.Set(strconv.Itoa(i), i)
		}()
	}
	wg.Wait()

	for i := range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			expected := i

			value := concurrentMap.Get(strconv.Itoa(i))
			got, _ := value.(int)
			if got != expected {
				t.Errorf("Got %d, expected %d\n", got, expected)
			}
		}()
	}
	wg.Wait()
}
