package calendar

import (
	"context"
	"dictionaries-service/tenant"
	"fmt"
	"time"
)

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

type service struct {
	repository Repository
}

type Service interface {
	LoadByTypeAndRange(ctx context.Context, calendarType string, from, to time.Time) ([]DictionaryCalendar, error)
}

func (s *service) LoadByTypeAndRange(ctx context.Context, calendarType string, from, to time.Time) ([]DictionaryCalendar, error) {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't read tenant from context")
	}

	return s.repository.LoadByTypeAndRange(t.Name, calendarType, from, to)
}
