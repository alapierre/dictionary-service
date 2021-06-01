package http

import (
	"context"
	"dictionaries-service/calendar"
	commons "dictionaries-service/transport/http"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// swagger:parameters Calendar
type calendarRequest struct {
	CalendarType string
	DayFrom      time.Time
	DayTo        time.Time
}

//swagger:response saveCalendar
type calendarResponse struct {
	Day    string            `json:"day"`
	Tenant string            `json:"tenant,omitempty"`
	Name   *string           `json:"name"`
	Kind   *string           `json:"kind,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// swagger:parameters
type calendarDelete struct {
	Day          time.Time
	CalendarType string
}

func MakeLoadCalendarEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(calendarRequest)
		cal, err := service.LoadByTypeAndRange(ctx, req.CalendarType, req.DayFrom, req.DayTo)

		if err != nil {
			return commons.MakeRestError(err, "Can't load calendar type "+req.CalendarType)
		}

		var res []calendarResponse

		for _, c := range cal {
			res = append(res, calendarResponse{
				Day:    c.Day.Format("2006-01-02"),
				Tenant: c.Tenant,
				Name:   c.Name,
				Kind:   c.Kind,
				Labels: c.Labels,
			})
		}
		return res, nil
	}
}

func DecodeLoadCalendarRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)

	from, err := time.Parse("2006-01-02", vars["from"])
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("2006-01-02", vars["to"])
	if err != nil {
		return nil, err
	}

	return calendarRequest{
		CalendarType: vars["type"],
		DayFrom:      from,
		DayTo:        to,
	}, nil
}

func MakeSaveCalendarEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(calendar.SaveDto)
		err := service.Save(ctx, &req)

		if err != nil {
			return commons.MakeRestError(err, "cant_create_new_dictionary_entry")
		}

		return nil, nil
	}
}

func MakeUpdateCalendarEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(calendar.SaveDto)
		err := service.Update(ctx, &req)

		if err != nil {
			return commons.MakeRestError(err, "cant_create_new_dictionary_entry")
		}

		return nil, nil
	}
}

func DecodeSaveCalendarRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)

	day, err := time.Parse("2006-01-02", vars["day"])
	if err != nil {
		return nil, err
	}

	var request calendar.SaveDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.CalendarType = vars["type"]
	request.Day = day

	return request, nil
}

func MakeDeleteCalendar(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(calendarDelete)
		return nil, service.Delete(ctx, req.CalendarType, req.Day)
	}
}

func DecodeDeleteCalendarRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)

	day, err := time.Parse("2006-01-02", vars["day"])
	if err != nil {
		return nil, err
	}

	return calendarDelete{
		Day:          day,
		CalendarType: vars["type"],
	}, nil
}
