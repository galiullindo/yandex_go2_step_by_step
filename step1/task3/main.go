package main

import (
	"strings"
)

type UpperWriter struct {
	UpperString string
}

func (w *UpperWriter) Write(p []byte) (n int, err error) {
	w.UpperString = strings.ToUpper(string(p))
	return len(p), nil
}
