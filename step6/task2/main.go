package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	StudentNotFoundError = errors.New("student not found ")
	InternalServerError  = errors.New("internal server error")
	BadRequestError      = errors.New("bad request")
	ErrNamesIsEmpty      = errors.New("names is empty")
)

type Student struct {
	Name string
	Mark int
	Err  error
}

func fetch(name string) Student {
	var client http.Client
	student := Student{Name: name}

	url := fmt.Sprintf("http://localhost:8082/mark?name=%s", name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		student.Err = err
		return student
	}

	resp, err := client.Do(req)
	if err != nil {
		student.Err = err
		return student
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		student.Err = InternalServerError
		return student
	case http.StatusNotFound:
		student.Err = StudentNotFoundError
		return student
	case http.StatusBadRequest:
		student.Err = BadRequestError
		return student
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		student.Err = err
		return student
	}

	mark, err := strconv.Atoi(string(body))
	if err != nil {
		student.Err = err
		return student
	}

	student.Mark = mark
	return student
}

func Average(names []string) (int, error) {
	if len(names) == 0 {
		return 0, ErrNamesIsEmpty
	}

	ch := func(names []string) <-chan Student {
		ch := make(chan Student)
		go func(names []string) {
			defer close(ch)
			for _, name := range names {
				ch <- fetch(name)
			}
		}(names)
		return ch
	}(names)

	sum := 0
	for student := range ch {
		if student.Err != nil {
			return 0, student.Err
		}
		sum += student.Mark
	}

	avg := sum / len(names)

	return avg, nil
}
