package main

import (
	"context"
	"errors"
	"io"
)

var (
	ErrEmptySequence = errors.New("sequence is empty")
)

type Message struct {
	Bytes []byte
	N     int
	Err   error
}

func readMessage(reader io.Reader, p []byte) Message {
	n, err := reader.Read(p)
	if err != nil {
		if err == io.EOF {
			return Message{p, n, io.EOF}
		} else {
			return Message{nil, 0, err}
		}
	}
	return Message{p, n, err}
}

func readWithContext(ctx context.Context, reader io.Reader) <-chan Message {
	channel := make(chan Message)

	go func() {
		defer close(channel)

		bytes := make([]byte, 16)
		for {
			select {
			case <-ctx.Done():
				return
			case channel <- readMessage(reader, bytes):
			}
		}
	}()

	return channel
}

func checkingForSequence(message Message, sequence []byte, last int) (bool, int) {
	for i := 0; i < message.N; i++ {
		isEquals := message.Bytes[i] == sequence[last]
		if isEquals {
			last++
		} else {
			last = 0
		}

		isContainsSequence := last == len(sequence)
		if isContainsSequence {
			return true, last
		}
	}
	return false, last
}

func Contains(ctx context.Context, reader io.Reader, sequence []byte) (is bool, e error) {
	isEmptySequence := len(sequence) == 0
	if isEmptySequence {
		return false, ErrEmptySequence
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	reading := readWithContext(ctx, reader)

	var (
		sLastByte          int
		isContainsSequence bool
		hasEOF             bool
	)

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case message := <-reading:
			err := message.Err
			if err != nil {
				if err == io.EOF {
					hasEOF = true
				} else {
					return false, err
				}
			}

			if message.N > 0 {
				isContainsSequence, sLastByte = checkingForSequence(message, sequence, sLastByte)
			}

			if isContainsSequence || hasEOF {
				return isContainsSequence, nil
			}
		}
	}
}
