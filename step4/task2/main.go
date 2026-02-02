package main

import "sync"

type Count interface {
	Increment()
	GetValue() int
}

type Counter struct {
	value int
	mu    sync.RWMutex
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) GetValue() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}
