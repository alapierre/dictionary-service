package calendar

import (
	"context"
	"dictionaries-service/tenant"
	"fmt"
	"time"
)

func NewService(repository Repository, typeRepository TypeRepository) Service {
	return &service{repository: repository, typeRepository: typeRepository}
}

type service struct {
	repository     Repository
	typeRepository TypeRepository
}

type Service interface {
	LoadByTypeAndRange(ctx context.Context, calendarType string, from, to time.Time) ([]DictionaryCalendar, error)
	Save(ctx context.Context, calendar *SaveDto) error
	Update(ctx context.Context, calendar *SaveDto) error
	Delete(ctx context.Context, calendarType string, day time.Time) error
	LoadTypes(ctx context.Context) ([]DictionaryCalendarType, error)
}

func (s *service) Save(ctx context.Context, calendar *SaveDto) error {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return fmt.Errorf("can't read tenant from context")
	}

	return s.repository.Save(mapCalendar(calendar, t))
}

func (s *service) Update(ctx context.Context, calendar *SaveDto) error {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return fmt.Errorf("can't read tenant from context")
	}

	return s.repository.Update(mapCalendar(calendar, t))
}

func (s *service) SaveType(ctx context.Context, calendarType *DictionaryCalendarType) error {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return fmt.Errorf("can't read tenant from context")
	}

	calendarType.Tenant = t.Name
	return s.typeRepository.Save(calendarType)
}

func mapCalendar(calendar *SaveDto, t tenant.Tenant) *DictionaryCalendar {
	return &DictionaryCalendar{
		Day:    time.Time(calendar.Day),
		Tenant: t.Name,
		Name:   &calendar.Name,
		Type:   calendar.CalendarType,
		Kind:   calendar.Kind,
		Labels: calendar.Labels,
	}
}

func (s *service) Delete(ctx context.Context, calendarType string, day time.Time) error {
	t, ok := tenant.FromContext(ctx)
	if !ok {
		return fmt.Errorf("can't read tenant from context")
	}
	return s.repository.Delete(t.Name, calendarType, day)
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

func (s *service) LoadTypes(ctx context.Context) ([]DictionaryCalendarType, error) {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't read tenant from context")
	}

	return s.typeRepository.LoadAll(t.Name)
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
