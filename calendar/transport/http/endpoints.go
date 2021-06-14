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

func MakeLoadCalendarTypesEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		cal, err := service.LoadTypes(ctx)
		if err != nil {
			return commons.MakeRestError(err, "Can't load calendar types")
		}

		return cal, nil
	}
}

func MakeCreateCalendarTypesEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(calendar.DictionaryCalendarType)
		err := service.SaveType(ctx, &req)

		if err != nil {
			return commons.MakeRestError(err, "cant_create_new_dictionary_entry")
		}

		return nil, nil
	}
}

func MakeUpdateCalendarTypesEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(calendar.DictionaryCalendarType)
		err := service.UpdateType(ctx, &req)

		if err != nil {
			return commons.MakeRestError(err, "cant_create_new_dictionary_entry")
		}

		return nil, nil
	}
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
			return commons.MakeRestError(err, "cant_update_dictionary_entry")
		}

		return nil, nil
	}
}

func DecodeSaveCalendarRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)

	var request calendar.SaveDto
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	request.CalendarType = vars["type"]

	return request, nil
}

func DecodeSaveCalendarTypeRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var request calendar.DictionaryCalendarType
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
}

func MakeDeleteCalendar(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(calendarDelete)

		if err := service.Delete(ctx, req.CalendarType, req.Day); err != nil {
			return commons.MakeRestError(err, "cant_delete_dictionary_entry")
		}

		return nil, nil
	}
}

func MakeDeleteCalendarType(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(calendarTypeDeleteWrapper)

		if err := service.DeleteType(ctx, req.Type); err != nil {
			return commons.MakeRestError(err, "cant_delete_calendar_type")
		}

		return nil, nil
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

func DecodeDeleteCalendarTypeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return calendarTypeDeleteWrapper{
		Type: vars["type"],
	}, nil
}
