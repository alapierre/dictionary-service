package main

import (
	"dictionaries-service/calendar"
	rest "dictionaries-service/calendar/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeCalendarEndpoints(r *mux.Router, service calendar.Service) {

	// swagger:route GET /api/calendar/{calendarType}/{from}/{to} Calendar
	//
	// Loads all calendar items for given type and date range
	r.Methods("GET").Path("/api/calendar/{type}/{from}/{to}").Handler(http.NewServer(
		rest.MakeLoadCalendarEndpoint(service),
		rest.DecodeLoadCalendarRequest,
		common.EncodeResponse,
	))

	// swagger:route POST /api/calendar/{calendarType}/{day} saveCalendar
	//
	// Create new calendar item
	r.Methods("POST").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeSaveCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	// swagger:route PUT /api/calendar/{calendarType}/{day} updateCalendar
	//
	// Update existing calendar item
	r.Methods("PUT").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeUpdateCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	// swagger:route DELETE /api/calendar/{calendarType}/{day} deleteCalendar
	//
	// Delete existing calendar item
	r.Methods("DELETE").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeDeleteCalendar(service),
		rest.DecodeDeleteCalendarRequest,
		common.EncodeSavedResponse,
	))
}
