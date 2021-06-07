package main

import (
	"dictionaries-service/confguration"
	http2 "dictionaries-service/confguration/transport/http"
	rest "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeConfigurationEndpoints(r *mux.Router, configurationService confguration.Service) {

	// swagger:route GET /api/config/{key}/{day} loadConfig
	//
	// Load requested config item in given day
	// Responses:
	//   200: loadConfigurationResponse
	r.Methods("GET").Path("/api/config/{key}/{day}").Handler(http.NewServer(
		http2.MakeLoadConfigurationEndpoint(configurationService),
		http2.DecodeLoadConfigurationRequest,
		rest.EncodeResponse,
	))

	// swagger:route GET /api/config/{day} loadAllConfigForDay
	//
	// Load config item for requested keys in given day
	// Responses:
	//   200: loadConfigurationResponse
	r.Methods("GET").Path("/api/configs/{day}").Handler(http.NewServer(
		http2.MakeLoadConfigurationArrayEndpoint(configurationService),
		http2.DecodeLoadConfigurationArrayRequest,
		rest.EncodeResponse,
	))
}
