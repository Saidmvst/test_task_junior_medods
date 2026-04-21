package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	dates := s.calculateDates(normalized.Recurrence)
	now := s.now()

	if len(dates) == 0 {
		model := &taskdomain.Task{
			Title:       normalized.Title,
			Description: normalized.Description,
			Status:      normalized.Status,
			CreatedAt:   now,
			UpdatedAt:   now,
			Recurrence:  normalized.Recurrence,
		}
		return s.repo.Create(ctx, model)
	}

	var firstCreated *taskdomain.Task
	for i, date := range dates {
		model := &taskdomain.Task{
			Title:       normalized.Title,
			Description: normalized.Description,
			Status:      normalized.Status,
			CreatedAt:   date,
			UpdatedAt:   now,
			Recurrence:  normalized.Recurrence,
		}
		created, err := s.repo.Create(ctx, model)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			firstCreated = created
		}
	}

	return firstCreated, nil
}

func (s *Service) calculateDates(settings *taskdomain.RecurrenceSettings) []time.Time {
	if settings == nil {
		return nil
	}

	var dates []time.Time
	start := s.now()
	end := start.AddDate(0, 0, 30)

	switch settings.Type {
	case taskdomain.Daily:
		for d := start; d.Before(end); d = d.AddDate(0, 0, settings.Interval) {
			dates = append(dates, d)
		}
	case taskdomain.Parity:
		for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
			isEven := d.Day()%2 == 0
			if (settings.Parity == "even" && isEven) || (settings.Parity == "odd" && !isEven) {
				dates = append(dates, d)
			}
		}
	case taskdomain.Monthly:
		d := time.Date(start.Year(), start.Month(), settings.DayOfMonth, start.Hour(), start.Minute(), start.Second(), 0, start.Location())
		if d.Before(start) {
			d = d.AddDate(0, 1, 0)
		}
		for d.Before(end) {
			dates = append(dates, d)
			d = d.AddDate(0, 1, 0)
		}
	case taskdomain.SpecificDays:
		return settings.SpecificDates
	}

	return dates
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}
	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}
	return s.repo.Update(ctx, model)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}
	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}
	if input.Recurrence != nil {
		switch input.Recurrence.Type {
		case taskdomain.Daily:
			if input.Recurrence.Interval <= 0 {
				return CreateInput{}, fmt.Errorf("%w: interval must be positive", ErrInvalidInput)
			}
		case taskdomain.Monthly:
			if input.Recurrence.DayOfMonth < 1 || input.Recurrence.DayOfMonth > 31 {
				return CreateInput{}, fmt.Errorf("%w: invalid day of month", ErrInvalidInput)
			}
		}
	}
	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}
	return input, nil
}
