package task

import (
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	Daily        RecurrenceType = "daily"
	Monthly      RecurrenceType = "monthly"
	SpecificDays RecurrenceType = "specific_days"
	Parity       RecurrenceType = "parity"
)

type RecurrenceSettings struct {
	Type          RecurrenceType `json:"type"`
	Interval      int            `json:"interval,omitempty"`
	DayOfMonth    int            `json:"day_of_month,omitempty"`
	SpecificDates []time.Time    `json:"specific_dates,omitempty"`
	Parity        string         `json:"parity,omitempty"`
}

type Task struct {
	ID          int64               `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      Status              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Recurrence  *RecurrenceSettings `json:"recurrence,omitempty"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
