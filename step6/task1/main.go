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
	NotFoundError       = errors.New("student not found ")
	InternalServerError = errors.New("internal server error")
	BadRequestError     = errors.New("bad request")
	NamesIsEmptyError   = errors.New("names is empty")
)

type Student struct {
	Name string
	Mark int
}

type Message struct {
	Student Student
	Err     error
}

func fetchMark(name string) (Student, error) {
	var client http.Client
	var student Student

	url := fmt.Sprintf("http://localhost:8082/mark?name=%s", name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return student, err
	}

	student.Name = name

	resp, err := client.Do(req)
	if err != nil {
		return student, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return student, BadRequestError
	}
	if resp.StatusCode == http.StatusNotFound {
		return student, NotFoundError
	}
	if resp.StatusCode == http.StatusInternalServerError {
		return student, InternalServerError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return student, err
	}

	mark, err := strconv.Atoi(string(body))
	if err != nil {
		return student, err
	}

	student.Mark = mark
	return student, nil
}

func Compare(name1 string, name2 string) (string, error) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	var message1, message2 Message

	wg.Go(func() {
		student, err := fetchMark(name1)
		mu.Lock()
		defer mu.Unlock()
		message1 = Message{Student: student, Err: err}
	})

	wg.Go(func() {
		student, err := fetchMark(name2)
		mu.Lock()
		defer mu.Unlock()
		message2 = Message{Student: student, Err: err}
	})

	wg.Wait()

	if message1.Err != nil {
		return "", message1.Err
	}
	if message2.Err != nil {
		return "", message2.Err
	}

	if message1.Student.Mark > message2.Student.Mark {
		return ">", nil
	}
	if message1.Student.Mark < message2.Student.Mark {
		return "<", nil
	}

	return "=", nil
}
