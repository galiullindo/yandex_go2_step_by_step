package testutils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var students = map[string]int{
	"Sara60":             60,
	"Bob50":              50,
	"Jack50":             50,
	"John40":             40,
	"Den10":              10,
	"Barbara25&ise=true": 0, // for returning internal server error status.
}

type Params struct {
	Name string
	ISE  bool
}

func ParseParams(r *http.Request) (Params, error) {
	query := r.URL.Query()
	params := Params{}

	params.Name = query.Get("name")
	if params.Name == "" {
		return params, fmt.Errorf("missing name")
	}

	iseStr := query.Get("ise")
	if iseStr == "" {
		params.ISE = false
	} else {
		ise, err := strconv.ParseBool(iseStr)
		if err != nil {
			params.ISE = false
		} else {
			params.ISE = ise
		}
	}

	return params, nil
}

func Mark(w http.ResponseWriter, r *http.Request) {
	params, err := ParseParams(r)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if params.ISE {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	mark, found := students[params.Name]
	if !found {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%d", mark)
}

func NewServer(addr string) (server *http.Server, start func(), stop func()) {
	mux := http.NewServeMux()
	mux.HandleFunc("/mark", Mark)

	server = &http.Server{Addr: addr, Handler: mux}

	start = func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("error in the server: %s\n", err)
		}
	}

	stop = func() {
		if err := server.Close(); err != nil {
			log.Printf("error stopping the server: %s\n", err)
		}
	}

	return server, start, stop
}
