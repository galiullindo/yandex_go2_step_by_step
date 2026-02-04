package main

import (
	"context"
	"errors"
	"io"
)

var (
	ErrEmptySequence = errors.New("sequence is empty")
)

type message struct {
	Bytes []byte
	N     int
	Err   error
}

func read(reader io.Reader, channel chan message) {
	defer close(channel)
	bytes := make([]byte, 16)

	for {
		n, err := reader.Read(bytes)

		if err == io.EOF {
			channel <- message{bytes, n, nil}
			return
		} else if err != nil {
			channel <- message{nil, 0, err}
			return
		}

		channel <- message{bytes, n, err}
	}
}

func Contains(ctx context.Context, reader io.Reader, sequence []byte) (bool, error) {
	isEmptySequence := len(sequence) == 0
	if isEmptySequence {
		return false, ErrEmptySequence
	}

	channel := make(chan message)

	go read(reader, channel)

	sLastByte := 0
	isContainsSequence := false
	for {
		select {
		case message, ok := <-channel:
			if !ok {
				return isContainsSequence, nil
			}

			err := message.Err
			if err != nil {
				return false, err
			}

			if message.N > 0 {
				for i := 0; i < message.N; i++ {
					isEquals := message.Bytes[i] == sequence[sLastByte]
					if isEquals {
						sLastByte++
					} else {
						sLastByte = 0
					}

					isContainsSequence = sLastByte == len(sequence)
					if isContainsSequence {
						return true, nil
					}
				}
			}
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}
}
