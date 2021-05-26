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

	if t.MergeDefault {
		return s.mergeByTypeAndRange(t.Name, calendarType, from, to)
	}

	return s.repository.LoadByTypeAndRange(t.Name, calendarType, from, to)
}

func (s *service) mergeByTypeAndRange(tenant, calendarType string, from, to time.Time) ([]DictionaryCalendar, error) {

	buf := make(map[time.Time]DictionaryCalendar)

	formDefault, err := s.repository.LoadByTypeAndRange("", calendarType, from, to)
	if err != nil {
		return nil, err
	}

	updateMap(buf, formDefault)

	formTenant, err := s.repository.LoadByTypeAndRange(tenant, calendarType, from, to)
	if err != nil {
		return nil, err
	}

	updateMap(buf, formTenant)

	return mapToSlice(buf), nil
}

func updateMap(buf map[time.Time]DictionaryCalendar, calendar []DictionaryCalendar) {
	for _, c := range calendar {
		buf[c.Day] = c
	}
}

func mapToSlice(buf map[time.Time]DictionaryCalendar) (result []DictionaryCalendar) {
	for _, v := range buf {
		result = append(result, v)
	}
	return
}
