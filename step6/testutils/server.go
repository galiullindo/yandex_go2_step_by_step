package testutils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var students = map[string]int{
	"Sara60":               60,
	"Bob50":                50,
	"Jack50":               50,
	"John40":               40,
	"Den10":                10,
	"Barbara25&ise=true":   0, // for returning internal server error status.
	"Barbara25&abort=true": 0, // for panic with error abort handler.
	"Barbara25&read=true":  0, // for error in read body.
	"Barbara25&conv=true":  0, // for returnig strng instead of integer.
}

type Params struct {
	Name  string
	ISE   bool
	Abort bool
	Read  bool
	Conv  bool
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

	abortStr := query.Get("abort")
	if abortStr == "" {
		params.Abort = false
	} else {
		abort, err := strconv.ParseBool(abortStr)
		if err != nil {
			params.Abort = false
		} else {
			params.Abort = abort
		}
	}

	readStr := query.Get("read")
	if readStr == "" {
		params.Read = false
	} else {
		read, err := strconv.ParseBool(readStr)
		if err != nil {
			params.Read = false
		} else {
			params.Read = read
		}
	}

	convStr := query.Get("conv")
	if convStr == "" {
		params.Conv = false
	} else {
		conv, err := strconv.ParseBool(convStr)
		if err != nil {
			params.Conv = false
		} else {
			params.Conv = conv
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

	if params.Abort {
		panic(http.ErrAbortHandler)
	}

	if params.Read {
		w.WriteHeader(http.StatusOK)
		hijacker, _ := w.(http.Hijacker)
		conn, bufrw, _ := hijacker.Hijack()
		_ = bufrw.Flush()
		conn.Close()
		return
	}

	if params.Conv {
		fmt.Fprintf(w, "%s", "*mark*")
		return
	}

	mark, found := students[params.Name]
	if !found {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%d", mark)
}

// NewServer возвращает указатель на экземпляр http.Server и функции start, и stop, для запуска и остановки сервера.
//
//	use:
//		_, start, stop := NewServer(addr)
//		go start()
//		defer stop()
//
// ...
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
