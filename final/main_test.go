package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

var diff = 5 * time.Millisecond

var ErrCustom = errors.New("custom error")

type CustomReader struct {
	f func(p []byte) (n int, err error)
}

func NewCustomReader(f func(p []byte) (n int, err error)) *CustomReader {
	return &CustomReader{f: f}
}

func (r *CustomReader) Read(b []byte) (n int, err error) {
	n, err = r.f(b)
	return n, err
}

func TestIsStatus(t *testing.T) {
	var tests = []struct {
		name     string
		s        string
		expected bool
	}{
		{name: "Case status Ready", s: "Готово", expected: true},
		{name: "Case status InProgress", s: "В работе", expected: true},
		{name: "Case status WillNotBeDone", s: "Не будет сделано", expected: true},
		{name: "Case random string", s: "абвгдеё", expected: false},
		{name: "Case empty string", s: "", expected: false},
		{name: "Case sensitivity check", s: "ГОТОВО", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := IsStatus(test.s); got != test.expected {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.s, got, test.expected)
			}
		})
	}
}

func TestIsTicket(t *testing.T) {
	var tests = []struct {
		name     string
		s        string
		expected bool
	}{
		{name: "Case valid ticket", s: "TICKET-...", expected: true},
		{name: "Case invalid ticket", s: "TICKET...", expected: false},
		{name: "Case random string", s: "абвгдеё", expected: false},
		{name: "Case empty string", s: "", expected: false},
		{name: "Case sensitivity check", s: "Ticket-", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := IsTicket(test.s); got != test.expected {
				t.Errorf("unexpected value for %v: got %v, expected %v\n", test.s, got, test.expected)
			}
		})
	}
}

func TestNewTicket(t *testing.T) {
	var tests = []struct {
		name        string
		ticket      string
		user        string
		status      string
		date        time.Time
		expected    *Ticket
		expectedErr error
	}{
		{
			name:     "Case valid ticket",
			ticket:   "TICKET-12345",
			user:     "user",
			status:   "Готово",
			date:     time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC),
			expected: &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
		},
		{
			name:        "Case invalid ticket's ticket",
			ticket:      "TICKET 12345",
			user:        "user",
			status:      "Готово",
			date:        time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC),
			expectedErr: ErrNotTicket,
		},
		{
			name:        "Case invalid ticket's status",
			ticket:      "TICKET-12345",
			user:        "user",
			status:      "Готово2",
			date:        time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC),
			expectedErr: ErrNotStatus,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewTicket(test.ticket, test.user, test.status, test.date)

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}

func TestParseTicket(t *testing.T) {
	var tests = []struct {
		name        string
		s           string
		expected    *Ticket
		expectedErr error
	}{
		{
			name:     "Case valid ticket",
			s:        "TICKET-12345_user_Готово_2026-01-03",
			expected: &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
		},
		{
			name:        "Case invalid ticket's ticket",
			s:           "TICKET 12345_user_Готово_2026-01-03",
			expectedErr: ErrNotTicket,
		},
		{
			name:        "Case invalid ticket's status",
			s:           "TICKET-12345_user_ВсёГотово_2026-01-03",
			expectedErr: ErrNotStatus,
		},
		{
			name:        "Case parsing error",
			s:           "TICKET-12345 user_Готово_2026-01-03",
			expectedErr: ErrParse,
		},
		{
			name:        "Case parsing date error",
			s:           "TICKET-12345_user_Готово_время...",
			expectedErr: ErrParse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseTicket(test.s, "_", "2006-01-02")

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}

func TestIsTarget(t *testing.T) {
	var (
		emptyUser  string = ""
		targetUser string = "user"
		anoterUser string = "another"

		emptyStatus  string = ""
		targetStatus string = "Готово"
		anoterStatus string = "Победа!"
	)

	var tests = []struct {
		name     string
		user     *string
		status   *string
		ticket   *Ticket
		expected bool
	}{
		{
			name:     "Case target ticket",
			user:     &targetUser,
			status:   &targetStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "Case any user",
			user:     nil,
			status:   &targetStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "Case any status",
			user:     &targetUser,
			status:   nil,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "Case any status and any user",
			user:     nil,
			status:   nil,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: true,
		},
		{
			name:     "Case empty user",
			user:     &emptyUser,
			status:   &targetStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "Case another user",
			user:     &anoterUser,
			status:   &targetStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "Case empty status",
			user:     &targetUser,
			status:   &emptyStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "Case another status",
			user:     &targetUser,
			status:   &anoterStatus,
			ticket:   &Ticket{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 3, 00, 0, 0, 0, time.UTC)},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.ticket.IsTarget(test.user, test.status)
			if got != test.expected {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}

func TestReadLines(t *testing.T) {
	var tests = []struct {
		name     string
		timeout  time.Duration
		reader   io.Reader
		expected []Line
	}{
		{
			name:     "Case normal",
			timeout:  10 * time.Millisecond,
			reader:   strings.NewReader("abc\ndef\n\nghi\njklmnop\nqrs\n"),
			expected: []Line{{"abc", nil}, {"def", nil}, {"", nil}, {"ghi", nil}, {"jklmnop", nil}, {"qrs", nil}},
		},
		{
			name:    "Case timeout 0ms",
			timeout: 0 * time.Millisecond,
			reader:  strings.NewReader("abc\ndef\n\nghi\njklmnop\nqrs\n"),
		},
		{
			name:    "Case delaytion",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(20 * time.Millisecond)
				n = copy(p, []byte("abcdefg"))
				return n, nil
			}),
		},
		{
			name:    "Case endless reading",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("abcdefg"))
				return n, nil
			}),
		},
		{
			name:    "Case read error",
			timeout: 10 * time.Millisecond,
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("abcdefg"))
				return n, ErrCustom
			}),
			expected: []Line{{"abcdefg", ErrCustom}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			start := time.Now()
			lines := ReadLines(ctx, test.reader)

			got := make([]Line, 0)
			for line := range lines {
				got = append(got, line)
			}

			duration := time.Since(start)
			timeout := test.timeout + diff
			if duration > timeout {
				t.Errorf("unexpected execution time: got %v, expected near %v\n", duration, timeout)
			}

			if !slices.Equal(got, test.expected) {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	var (
		user   string = "user"
		anoter string = "another"

		ready      string = "Готово"
		inProgress string = "В работе"
	)
	var lines = []string{
		"TICKET-12345_user_Готово_2026-01-02",
		"TICKET-12346_user_В работе_2026-01-03",
		"TICKET-12347_user_В работе_2026-01-04",
		"invalid-ticket_user_Готово_2026-01-03",       // invalid ticket
		"TICKET-12348_user_invalid-status_2026-01-03", // invalid status
		"TICKET-12349_user_Готово_invalid-data",       // invalid date
		"TICKET-12346_user_В_работе_2026-01-03",       // invalid separation 5
		"some-text",
		"TICKET-12355_another_Готово_2026-01-02",
		"TICKET-12356_another_В работе_2026-01-03",
		"TICKET-12357_another_Готово_2026-01-04",
		"invalid-ticket_another_Готово_2026-01-03",       // invalid ticket
		"TICKET-12358_another_invalid-status_2026-01-03", // invalid status
		"TICKET-12359_another_Готово_invalid-data",       // invalid date
		"TICKET-12346_another В работе_2026-01-03",       // invalid separation 3
		"se_pa_ra_ted-text",
	}

	var tests = []struct {
		name        string
		reader      io.Reader
		writer      io.Writer
		user        *string
		status      *string
		timeout     time.Duration
		expected    string
		expectedErr error
	}{
		{
			name:    "Case all tickets",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    nil,
			status:  nil,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12346", User: "user", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12347", User: "user", Status: "В работе", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12355", User: "another", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12356", User: "another", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12357", User: "another", Status: "Готово", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case only user's tickets",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    &user,
			status:  nil,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12346", User: "user", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12347", User: "user", Status: "В работе", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case only another's tickets",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    &anoter,
			status:  nil,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12355", User: "another", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12356", User: "another", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12357", User: "another", Status: "Готово", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case only ready tickets",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    nil,
			status:  &ready,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12345", User: "user", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12355", User: "another", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12357", User: "another", Status: "Готово", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case only in progress tickets",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    nil,
			status:  &inProgress,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12346", User: "user", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12347", User: "user", Status: "В работе", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12356", User: "another", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case user's tickets in progress",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    &user,
			status:  &inProgress,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12346", User: "user", Status: "В работе", Date: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12347", User: "user", Status: "В работе", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name:    "Case annother's tickets ready",
			reader:  strings.NewReader(strings.Join(lines, "\n")),
			writer:  bytes.NewBuffer(nil),
			user:    &anoter,
			status:  &ready,
			timeout: 10 * time.Millisecond,
			expected: func() string {
				s, _ := json.Marshal([]Ticket{
					{Ticket: "TICKET-12355", User: "another", Status: "Готово", Date: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Ticket: "TICKET-12357", User: "another", Status: "Готово", Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)},
				})
				return string(s)
			}(),
		},
		{
			name: "Case read error",
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("TICKET-12345_user_Готово_2026-01-02\n"))
				err = ErrCustom
				return n, err
			}),
			writer:      bytes.NewBuffer(nil),
			timeout:     10 * time.Millisecond,
			expectedErr: ErrCustom,
		},
		{
			name: "Case delation",
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				time.Sleep(20 * time.Millisecond)
				n = copy(p, []byte("TICKET-12345_user_Готово_2026-01-02\n"))
				return n, err
			}),
			writer:      bytes.NewBuffer(nil),
			timeout:     10 * time.Millisecond,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name: "Case endless reading",
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("TICKET-12345_user_Готово_2026-01-02\n"))
				return n, err
			}),
			writer:      bytes.NewBuffer(nil),
			timeout:     10 * time.Millisecond,
			expectedErr: context.DeadlineExceeded,
		},
		{
			name: "Case timeout 0",
			reader: NewCustomReader(func(p []byte) (n int, err error) {
				n = copy(p, []byte("TICKET-12345_user_Готово_2026-01-02\n"))
				return n, err
			}),
			writer:      bytes.NewBuffer(nil),
			timeout:     0 * time.Millisecond,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			start := time.Now()
			err := GetTasks(ctx, test.reader, test.writer, test.user, test.status, test.timeout)
			duration := time.Since(start)

			timeout := test.timeout + diff
			if duration > timeout {
				t.Errorf("unexpected execution time: got %v, expected near %v\n", duration, timeout)
			}

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("unexpected error: got %v, expected %v\n", err, test.expectedErr)
			}

			var got string
			if w, ok := test.writer.(*bytes.Buffer); ok {
				got = w.String()
			} else {
				t.Fatalf("unexpected writer: %v\n", test.writer)
			}

			if got != test.expected {
				t.Errorf("unexpected value: got %v, expected %v\n", got, test.expected)
			}
		})
	}
}
