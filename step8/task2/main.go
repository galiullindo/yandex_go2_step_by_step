package main

import (
	"context"
	"io"
	"os"
)

func ReadWithContext(ctx context.Context, r io.Reader, p []byte) (int, error) {
	type res struct {
		n   int
		err error
	}

	resCh := make(chan res)
	go func() {
		n, err := r.Read(p)
		resCh <- res{n, err}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case msg := <-resCh:
		return msg.n, msg.err
	}
}

func ReadAllWithContext(ctx context.Context, r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		select {
		case <-ctx.Done():
			return b, ctx.Err()
		default:
			n, err := ReadWithContext(ctx, r, b[len(b):cap(b)])
			b = b[:len(b)+n]
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				return b, err
			}

			if len(b) == cap(b) {
				b = append(b, 0)[:len(b)]
			}
		}
	}
}

func MakeChannelForReading(ctx context.Context, filePath string) <-chan []byte {
	channel := make(chan []byte)

	go func() {
		defer close(channel)

		select {
		case <-ctx.Done():
			channel <- nil
			return
		default:
			file, err := os.Open(filePath)
			if err != nil {
				channel <- nil
				return
			}
			defer file.Close()

			bytes, err := ReadAllWithContext(ctx, file)
			if err != nil {
				channel <- nil
				return
			}

			channel <- bytes
		}
	}()

	return channel
}

func ReadJSON(ctx context.Context, filePath string, forSend chan<- []byte) {
	defer close(forSend)
	var forReceive <-chan []byte

	select {
	case <-ctx.Done():
		return
	default:
		forReceive = MakeChannelForReading(ctx, filePath)
	}

	select {
	case <-ctx.Done():
		return
	case bytes := <-forReceive:
		if bytes != nil {
			forSend <- bytes
		}
		return
	}
}
