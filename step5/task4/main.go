package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func StartServer(timeout time.Duration) {
	mux := http.NewServeMux()
	mux.HandleFunc("/readSource", func(w http.ResponseWriter, r *http.Request) {
		var client http.Client

		request, err := http.NewRequest(http.MethodGet, "http://localhost:8081/provideData", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		responce, err := client.Do(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		defer responce.Body.Close()

		io.Copy(w, responce.Body)
	})

	handlerWithTimeout := http.TimeoutHandler(mux, timeout, "timeout")

	server = &http.Server{Addr: ":8080", Handler: handlerWithTimeout}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("error in the server: %s\n", err)
		}
	}()
}
