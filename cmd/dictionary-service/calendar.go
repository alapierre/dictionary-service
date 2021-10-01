package main

import (
	"dictionaries-service/calendar"
	rest "dictionaries-service/calendar/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	nethttp "net/http"
)

func makeCalendarEndpoints(r *mux.Router, service calendar.Service) {

	// swagger:route GET /api/calendar loadCalendarTypes
	//
	// Loads all calendar types
	//     Responses:
	//       200: calendarTypeResponse
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/calendar").Handler(http.NewServer(
		rest.MakeLoadCalendarTypesEndpoint(service),
		common.EmptyRequest(),
		common.EncodeResponse,
	))

	// swagger:route POST /api/calendar saveCalendarType
	//
	// Create new calendar type
	//     Responses:
	//       201:
	//       400: RestError
	r.Methods("POST", "OPTIONS").Path("/api/calendar").Handler(http.NewServer(
		rest.MakeCreateCalendarTypesEndpoint(service),
		rest.DecodeSaveCalendarTypeRequest,
		common.EncodeWithStatus(201),
	))

	// swagger:route DELETE /api/calendar/{type} deleteCalendarType
	//
	// DeleteValue calendar type - will work only if no any calendar items for this type exist
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("DELETE", "OPTIONS").Path("/api/calendar/{type}").Handler(http.NewServer(
		rest.MakeDeleteCalendarType(service),
		rest.DecodeDeleteCalendarTypeRequest,
		common.EncodeWithStatus(204),
	))

	// swagger:route PUT /api/calendar updateCalendarType
	//
	// Create new calendar type
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("PUT", "OPTIONS").Path("/api/calendar").Handler(http.NewServer(
		rest.MakeUpdateCalendarTypesEndpoint(service),
		rest.DecodeSaveCalendarTypeRequest,
		common.EncodeWithStatus(204),
	))

	// swagger:route GET /api/calendar/{calendarType}/{from}/{to} loadCalendar
	//
	// Loads all calendar items for given type and date range
	//     Responses:
	//       200: calendarResponse
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/calendar/{calendarType}/{dayFrom}/{dayTo}").Handler(http.NewServer(
		rest.MakeLoadCalendarEndpoint(service),
		rest.DecodeLoadCalendarRequest,
		common.EncodeResponse,
	))

	// swagger:route POST /api/calendar/{type} saveCalendar
	//
	// Create new calendar item for given type and day
	//     Responses:
	//       201:
	//       400: RestError
	r.Methods("POST", "OPTIONS").Path("/api/calendar/{type}").Handler(http.NewServer(
		rest.MakeSaveCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeWithStatus(nethttp.StatusCreated),
	))

	// swagger:route PUT /api/calendar/{type} updateCalendar
	//
	// Update existing calendar item for given type and day
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("PUT", "OPTIONS").Path("/api/calendar/{type}").Handler(http.NewServer(
		rest.MakeUpdateCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	// swagger:route DELETE /api/calendar/{type}/{day} deleteCalendar
	//
	// DeleteValue existing calendar item for given type and day
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("DELETE", "OPTIONS").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeDeleteCalendar(service),
		rest.DecodeDeleteCalendarRequest,
		common.EncodeWithStatus(204),
	))
}
