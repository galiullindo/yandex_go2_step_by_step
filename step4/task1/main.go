package main

import "sync"

type SafeMap struct {
	m     map[string]interface{}
	mutex sync.Mutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{make(map[string]interface{}), sync.Mutex{}}
}

func (s *SafeMap) Get(key string) interface{} {
	s.mutex.Lock()
	value, found := s.m[key]
	s.mutex.Unlock()

	if !found {
		return nil
	}

	return value
}

func (s *SafeMap) Set(key string, value interface{}) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}
