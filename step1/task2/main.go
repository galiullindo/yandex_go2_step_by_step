package main

import (
	"io"
)

func ReadString(r io.Reader) (string, error) {
	stringBytes := make([]byte, 0)

	buf := make([]byte, 8)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			stringBytes = append(stringBytes, buf[:n]...)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return string(stringBytes), nil
}
