package main

import (
	"errors"
	"io"
)

var ErrSequenceLengthZero = errors.New("sequence length cannot be zero")

func checkingForSequence(buf []byte, seq []byte, n int, lt *int) bool {
	for i := 0; i < n; i++ {
		if buf[i] == seq[*lt] {
			*lt++
			if *lt == len(seq) {
				return true
			}
		} else {
			*lt = 0
		}
	}
	return false
}

func Contains(r io.Reader, seq []byte) (bool, error) {
	isContains := false
	lt := 0

	if len(seq) == 0 {
		return false, ErrSequenceLengthZero
	}

	buf := make([]byte, 16)
	for !isContains {
		n, err := r.Read(buf)
		if n > 0 {
			isContains = checkingForSequence(buf, seq, n, &lt)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
	}

	return isContains, nil
}
