package main

import (
	"dictionaries-service/calendar"
	rest "dictionaries-service/calendar/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeCalendarEndpoints(r *mux.Router, service calendar.Service) {

	// swagger:route GET /api/calendar/{calendarType}/{from}/{to} loadCalendar
	//
	// Loads all calendar items for given type and date range
	//     Responses:
	//       200: calendarResponse
	//       400: RestError
	r.Methods("GET").Path("/api/calendar/{calendarType}/{dayFrom}/{dayTo}").Handler(http.NewServer(
		rest.MakeLoadCalendarEndpoint(service),
		rest.DecodeLoadCalendarRequest,
		common.EncodeResponse,
	))

	// swagger:route POST /api/calendar/{type}/{day} saveCalendar
	//
	// Create new calendar item for given type and day
	//     Responses:
	//       201:
	//       400: RestError
	r.Methods("POST").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeSaveCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	// swagger:route PUT /api/calendar/{type}/{day} updateCalendar
	//
	// Update existing calendar item for given type and day
	//     Responses:
	//       200:
	//       400: RestError
	r.Methods("PUT").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeUpdateCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	// swagger:route DELETE /api/calendar/{type}/{day} deleteCalendar
	//
	// Delete existing calendar item for given type and day
	//     Responses:
	//       200:
	//       400: RestError
	r.Methods("DELETE").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeDeleteCalendar(service),
		rest.DecodeDeleteCalendarRequest,
		common.EncodeSavedResponse,
	))
}
