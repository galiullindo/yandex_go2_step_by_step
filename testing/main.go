package main

import (
	"errors"
	"strings"
	"time"
)

type TicketStatus string

const (
	Ready         TicketStatus = "Готово"
	InProgress    TicketStatus = "В работе"
	WillNotBeDone TicketStatus = "Не будет сделано"
)

var ErrInvalidStatus = errors.New("invalid status")

func IsStatus(status string) bool {
	return status == string(Ready) || status == string(InProgress) || status == string(WillNotBeDone)
}

type Ticket struct {
	Ticket string
	User   string
	Status string
	Date   time.Time
}

func NewTicket(ticket string, user string, status string, date time.Time) (Ticket, error) {
	if !IsStatus(status) {
		return Ticket{}, ErrInvalidStatus
	}
	return Ticket{Ticket: ticket, User: user, Status: status, Date: date}, nil

}

func isTicket(message string) bool {
	return strings.HasPrefix(message, "TICKET")
}

var ErrInvalidData = errors.New("invalid data")

func parseMessageToTicket(message string) (Ticket, error) {
	data := strings.Split(message, "_")
	if len(data) != 4 {
		return Ticket{}, ErrInvalidData
	}

	ticket := data[0]
	user := data[1]
	status := data[2]

	date, err := time.Parse("2006-01-02", data[3])
	if err != nil {
		return Ticket{}, err
	}

	newticket, err := NewTicket(ticket, user, status, date)
	if err != nil {
		return Ticket{}, err
	}

	return newticket, err
}

func checkTicketParams(ticket *Ticket, user *string, status *string) bool {
	if user != nil && *user != ticket.User {
		return false
	}
	if status != nil && *status != ticket.Status {
		return false
	}

	return true
}

func GetTasks(text string, user *string, status *string) []Ticket {
	tickets := []Ticket{}

	messages := strings.Split(text, "\n")
	for _, message := range messages {
		message = strings.TrimSpace(message)
		if !isTicket(message) {
			continue
		}

		ticket, err := parseMessageToTicket(message)
		if err != nil {
			continue
		}

		if checkTicketParams(&ticket, user, status) {
			tickets = append(tickets, ticket)
		}
	}

	return tickets
}
