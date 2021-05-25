package http

import (
	"context"
	"dictionaries-service/calendar"
	"dictionaries-service/transport/http"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	nethttp "net/http"
	"time"
)

type calendarRequest struct {
	CalendarType string
	DayFrom      time.Time
	DayTo        time.Time
}

type calendarResponse struct {
	Day    string            `json:"day"`
	Tenant string            `json:"tenant"`
	Name   *string           `json:"name"`
	Kind   *string           `json:"kind,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func MakeLoadCalendarEndpoint(service calendar.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(calendarRequest)
		cal, err := service.LoadByTypeAndRange(ctx, req.CalendarType, req.DayFrom, req.DayTo)

		if err != nil {
			return http.MakeRestError(err, "Can't load calendar type "+req.CalendarType)
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

func DecodeLoadConfigurationArrayRequest(_ context.Context, r *nethttp.Request) (interface{}, error) {

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
