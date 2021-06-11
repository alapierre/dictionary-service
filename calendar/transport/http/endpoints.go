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

// swagger:parameters loadCalendar
type calendarRequest struct {

	// in:path
	CalendarType string `json:"calendarType"`

	// swagger:strfmt date
	// in:path
	DayFrom time.Time `json:"from"`

	// swagger:strfmt date
	// in:path
	DayTo time.Time `json:"to"`

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// should combine response with given and default tenant
	// in:header
	Merge string `json:"X-Tenant-Merge-Default"`
}

type calendarResponse struct {
	Day    string            `json:"day"`
	Tenant string            `json:"tenant,omitempty"`
	Name   *string           `json:"name"`
	Kind   *string           `json:"kind,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

//swagger:response calendarResponse
type calendarResponseWrapper struct {
	// in:body
	Body calendarResponse
}

// swagger:parameters deleteCalendar
type calendarDelete struct {

	// swagger:strfmt date
	// in:path
	Day time.Time `json:"day"`

	// in:path
	CalendarType string `json:"type"`

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`
}

// swagger:parameters saveCalendar updateCalendar
type SaveDtoWrapper struct {

	// Calendar type id
	// in:path
	Type string

	// Day in calendar
	// swagger:strfmt date
	// in:path
	Day time.Time

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// should combine response with given and default tenant
	// in:header
	Merge string `json:"X-Tenant-Merge-Default"`

	// Calendar body
	// in:body
	Body calendar.SaveDto
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

	from, err := time.Parse("2006-01-02", vars["dayFrom"])
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("2006-01-02", vars["dayTo"])
	if err != nil {
		return nil, err
	}

	return calendarRequest{
		CalendarType: vars["calendarType"],
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
