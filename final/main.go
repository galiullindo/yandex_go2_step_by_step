package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

var (
	ErrNotStatus = errors.New("is not a status")
	ErrNotTicket = errors.New("is not a ticket")
	ErrParse     = errors.New("cannot be parsed to a ticket")
)

type Status string

const (
	Ready         Status = "Готово"
	InProgress    Status = "В работе"
	WillNotBeDone Status = "Не будет сделано"
)

func IsStatus(s string) bool {
	status := Status(s)
	return status == Ready || status == InProgress || status == WillNotBeDone
}

func IsTicket(s string) bool {
	return strings.HasPrefix(s, "TICKET-")
}

type Ticket struct {
	Ticket string
	User   string
	Status string
	Date   time.Time
}

func NewTicket(ticket string, user string, status string, date time.Time) (*Ticket, error) {
	if !IsTicket(ticket) {
		return nil, fmt.Errorf("%v %w", ticket, ErrNotTicket)
	}
	if !IsStatus(status) {
		return nil, fmt.Errorf("%v %w", status, ErrNotStatus)
	}
	return &Ticket{Ticket: ticket, User: user, Status: status, Date: date}, nil
}

func ParseTicket(s string, sep string, layout string) (*Ticket, error) {
	p := strings.Split(s, sep)
	if l := len(p); l != 4 {
		return nil, fmt.Errorf("%d piece %w", l, ErrParse)
	}

	ticket := p[0]
	user := p[1]
	status := p[2]

	date, err := time.Parse(layout, p[3])
	if err != nil {
		return nil, fmt.Errorf("%v %w", p[3], ErrParse)
	}

	return NewTicket(ticket, user, status, date)
}

func (t *Ticket) IsTarget(user *string, status *string) bool {
	if user != nil && *user != t.User {
		return false
	}
	if status != nil && *status != t.Status {
		return false
	}

	return true
}

type Line struct {
	Text string
	Err  error
}

func ReadLines(ctx context.Context, r io.Reader) <-chan Line {
	channel := make(chan Line)

	go func() {
		defer close(channel)

		ch := make(chan bool)
		scanner := bufio.NewScanner(r)
		for {
			go func() { ch <- scanner.Scan() }()

			select {
			case <-ctx.Done():
				return
			case scan := <-ch:
				if !scan {
					return
				}

				channel <- Line{Text: scanner.Text(), Err: scanner.Err()}
				if scanner.Err() != nil {
					return
				}
			}
		}
	}()

	return channel
}

func GetTasks(ctx context.Context, r io.Reader, w io.Writer, user *string, status *string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tikets := make([]Ticket, 0)
	lines := ReadLines(ctx, r)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-lines:
			if !ok {
				b, err := json.Marshal(tikets)
				if err != nil {
					return err
				}
				if _, err := w.Write(b); err != nil {
					return err
				}
				return nil
			}

			if line.Err != nil {
				return line.Err
			}

			theTicket, err := ParseTicket(strings.TrimSpace(line.Text), "_", "2006-01-02")
			if err != nil {
				continue
			}

			if theTicket.IsTarget(user, status) {
				tikets = append(tikets, *theTicket)
			}
		}
	}
}
