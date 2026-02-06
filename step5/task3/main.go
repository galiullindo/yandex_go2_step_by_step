package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

type APIResponse struct {
	URL        string
	Data       string
	StatusCode int
	Err        error
}

func Fetch(ctx context.Context, url string) APIResponse {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return APIResponse{URL: url, Err: err}
	}

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		if response == nil {
			return APIResponse{URL: url, Err: err}
		}
		return APIResponse{URL: url, StatusCode: response.StatusCode, Err: err}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return APIResponse{URL: url, StatusCode: response.StatusCode, Err: err}
	}

	return APIResponse{URL: url, Data: string(body), StatusCode: response.StatusCode, Err: err}
}

func FetchAPI(ctx context.Context, urls []string, timeout time.Duration) []APIResponse {
	responses := make([]APIResponse, 0, len(urls))
	channel := make(chan APIResponse)
	defer close(channel)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, url := range urls {
		go func() { channel <- Fetch(ctx, url) }()
	}

	for range urls {
		responses = append(responses, <-channel)
	}

	return responses
}
