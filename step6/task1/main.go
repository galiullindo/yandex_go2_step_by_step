package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var (
	StudentNotFoundError = errors.New("student not found ")
	InternalServerError  = errors.New("internal server error")
	BadRequestError      = errors.New("bad request")
)

func fetch(name string) (int, error) {
	var client http.Client
	url := fmt.Sprintf("http://localhost:8082/mark?name=%s", name)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return 0, InternalServerError
	case http.StatusNotFound:
		return 0, StudentNotFoundError
	case http.StatusBadRequest:
		return 0, BadRequestError
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	mark, err := strconv.Atoi(string(body))
	if err != nil {
		return 0, err
	}

	return mark, nil
}

type Student struct {
	Name string
	Mark int
	Err  error
}

func Compare(name1 string, name2 string) (string, error) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	var student1, student2 Student

	wg.Go(func() {
		mark, err := fetch(name1)
		mu.Lock()
		defer mu.Unlock()
		student1 = Student{Name: name1, Mark: mark, Err: err}
	})

	wg.Go(func() {
		mark, err := fetch(name2)
		mu.Lock()
		defer mu.Unlock()
		student2 = Student{Name: name2, Mark: mark, Err: err}
	})

	wg.Wait()

	if student1.Err != nil {
		return "", student1.Err
	}
	if student2.Err != nil {
		return "", student2.Err
	}

	if student1.Mark > student2.Mark {
		return ">", nil
	}
	if student1.Mark < student2.Mark {
		return "<", nil
	}

	return "=", nil
}
