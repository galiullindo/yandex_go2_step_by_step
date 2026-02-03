package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestConcurrenceQueue(t *testing.T) {
	var tests = []struct {
		name      string
		elements  []any
		expecteds []any
	}{
		{
			name:      "Enqueue and Dequeue case integer",
			elements:  []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expecteds: []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:      "Enqueue and Dequeue case string",
			elements:  []any{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			expecteds: []any{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			name:      "Enqueue and Dequeue case mixed",
			elements:  []any{1, "a", 2, "b", 3, "c", 4, "d", 5, "e", 6, "f", 7, "g", 8, "h", 9, "i", 0, "j"},
			expecteds: []any{1, "a", 2, "b", 3, "c", 4, "d", 5, "e", 6, "f", 7, "g", 8, "h", 9, "i", 0, "j"},
		},
		{
			name:      "Enqueue and Dequeue case empty queue  ",
			elements:  []any(nil),
			expecteds: []any{nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var queue Queue = &ConcurrentQueue{}

			for _, element := range test.elements {
				queue.Enqueue(element)
			}

			for _, expected := range test.expecteds {
				got := queue.Dequeue()
				if got != expected {
					t.Errorf("Got %v, expected %v\n", got, expected)
				}
			}
		})
	}
}

func TestConcurrenceQueueWithGoroutines(t *testing.T) {
	var tests = []struct {
		name       string
		goroutines int
		expected   int
	}{
		{
			name:       "Parallel processing case 100 goroutines and value 1",
			goroutines: 100,
			expected:   100,
		},
		{
			name:       "Parallel processing case 1000 goroutines and value 1",
			goroutines: 1000,
			expected:   1000,
		},
		{
			name:       "Parallel processing case 10000 goroutines and value 1",
			goroutines: 10000,
			expected:   10000,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var queue Queue = &ConcurrentQueue{}
			var wg sync.WaitGroup

			// тут значения тоже изменяются тк, указатель на то же значение что и queue
			// просто другой тип указателя
			cQueue, ok := queue.(*ConcurrentQueue)
			if !ok {
				t.Fatalf("unexpected queue type: %v\n", reflect.TypeOf(queue))
			}

			for range test.goroutines {
				wg.Go(func() {
					queue.Enqueue(1)
				})
			}
			wg.Wait()

			if got := len(cQueue.queue); got != test.expected {
				t.Errorf("Queue length got %d, expected %d\n", got, test.expected)
			}

			for range test.goroutines {
				wg.Go(func() {
					queue.Dequeue()
				})
			}
			wg.Wait()

			if got := len(cQueue.queue); got != 0 {
				t.Errorf("Queue length got %d, expected 0\n", got)
			}
		})
	}
}
