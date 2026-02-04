package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrNilContext = errors.New("nil context")
)

type APIResponse struct {
	Data       string
	StatusCode int
}

func fetchAPI(ctx context.Context, url string, timeout time.Duration) (*APIResponse, error) {
	var client http.Client

	if ctx == nil {
		return nil, ErrNilContext
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponse{Data: string(body), StatusCode: response.StatusCode}, nil
}
