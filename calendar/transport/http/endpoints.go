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

	// in: path
	CalendarType string

	// in: path
	DayFrom time.Time

	// in: path
	DayTo time.Time
}

//swagger:response loadCalendar
type calendarResponse struct {

	// Day in YYYY-MM-dd format
	Day string `json:"day"`

	Tenant string `json:"tenant,omitempty"`

	// Day name in calendar
	// required: true
	Name *string `json:"name"`

	// Describe day kind in calendar
	Kind *string `json:"kind,omitempty"`

	// Optional extra labels
	Labels map[string]string `json:"labels,omitempty"`
}

// swagger:parameters deleteCalendar
type calendarDelete struct {
	Day          time.Time
	CalendarType string
}

// SaveDto model for create new Calendar item
// swagger1:parameters createCalendar
// in: body
type SaveDto struct {
	Day          time.Time `json:"-"`
	CalendarType string    `json:"-"`

	// Name of day in calendar
	Name string `json:"name"`

	// Kind of day, any string value useful for you eg. holiday type
	Kind *string `json:"kind,omitempty"`

	// extra labels describes day in calendar, eg. holiday properties
	Labels map[string]string `json:"labels,omitempty"`
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
		req := request.(SaveDto)
		err := service.Save(ctx, &req)

		if err != nil {
			return commons.MakeRestError(err, "cant_create_new_dictionary_entry")
		}

		return nil, nil
	}
}

func MakeUpdateCalendarEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SaveDto)
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

	var request SaveDto
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
