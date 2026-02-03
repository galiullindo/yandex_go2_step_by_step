package main

import "sync"

type Queue interface {
	Enqueue(element any)
	Dequeue() any
}

type ConcurrentQueue struct {
	queue []any
	mutex sync.Mutex
}

func (q *ConcurrentQueue) Enqueue(element any) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, element)
}

func (q *ConcurrentQueue) Dequeue() any {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.queue) > 0 {
		element := q.queue[0]
		q.queue = q.queue[1:]
		return element
	}

	return nil
}
