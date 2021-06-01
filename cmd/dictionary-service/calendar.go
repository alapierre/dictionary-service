package main

import (
	"dictionaries-service/calendar"
	rest "dictionaries-service/calendar/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeCalendarEndpoints(r *mux.Router, service calendar.Service) {

	r.Methods("GET").Path("/api/calendar/{type}/{from}/{to}").Handler(http.NewServer(
		rest.MakeLoadCalendarEndpoint(service),
		rest.DecodeLoadCalendarRequest,
		common.EncodeResponse,
	))

	r.Methods("POST").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeSaveCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	r.Methods("PUT").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeUpdateCalendarEndpoint(service),
		rest.DecodeSaveCalendarRequest,
		common.EncodeSavedResponse,
	))

	r.Methods("DELETE").Path("/api/calendar/{type}/{day}").Handler(http.NewServer(
		rest.MakeDeleteCalendar(service),
		rest.DecodeDeleteCalendarRequest,
		common.EncodeSavedResponse,
	))
}
