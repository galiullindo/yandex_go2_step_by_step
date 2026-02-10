package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func averageMark(students []Student) int {
	if len(students) == 0 {
		return 0
	}

	sum := 0
	for _, student := range students {
		sum += student.Mark
	}
	return sum / len(students)
}

func greaterThan(students []Student, mark int) []Student {
	if len(students) == 0 {
		return []Student{}
	}

	best := make([]Student, 0)
	for _, student := range students {
		if student.Mark > mark {
			best = append(best, student)
		}
	}
	return best
}

func BestStudents(names []string) (string, error) {
	if len(names) == 0 {
		return "", NamesIsEmptyError
	}

	channel := func() <-chan Message {
		channel := make(chan Message)
		go func() {
			defer close(channel)
			for _, name := range names {
				student, err := fetchMark(name)
				channel <- Message{Student: student, Err: err}
			}
		}()
		return channel
	}()

	students := make([]Student, 0, len(names))
	for message := range channel {
		student, err := message.Student, message.Err
		if err != nil {
			return "", err
		}
		students = append(students, student)
	}

	mark := averageMark(students)
	bestStudents := greaterThan(students, mark)

	bestStudentsNames := make([]string, 0, len(bestStudents))
	for _, student := range bestStudents {
		bestStudentsNames = append(bestStudentsNames, student.Name)
	}

	return strings.Join(bestStudentsNames, ","), nil
}
