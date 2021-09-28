package main

import (
	"dictionaries-service/configuration"
	rest "dictionaries-service/configuration/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeConfigurationEndpoints(r *mux.Router, configurationService configuration.Service) {

	// swagger:route GET /api/config/{key} loadValues
	//
	// Loads all config for given key
	//     Responses:
	//       200: loadValuesResponseWrapper
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/config/{key}").Handler(http.NewServer(
		rest.MakeLoadValuesEndpoint(configurationService),
		rest.DecodeLoadValuesRequest,
		common.EncodeResponse,
	))

	// swagger:route GET /api/config loadAllConfigKeys
	//
	// Loads all unique config keys, name and type
	//     Responses:
	//       200: loadShortResponseWrapper
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/config").Handler(http.NewServer(
		rest.MakeLoadAllShortEndpoint(configurationService),
		common.EmptyRequest(),
		common.EncodeResponse,
	))

	// swagger:route GET /api/config/{key}/{day} loadConfig
	//
	// Loads config value for given kay and day
	//     Responses:
	//       200: loadConfigurationOneResponseWrapper
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/config/{key}/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationEndpoint(configurationService),
		rest.DecodeLoadConfigurationRequest,
		common.EncodeResponse,
	))

	// swagger:route GET /api/configs/{day} loadConfigForKeys
	//
	// Loads all calendar items for given type and date range
	//     Responses:
	//       200: loadConfigurationResponse
	//       400: RestError
	r.Methods("GET", "OPTIONS").Path("/api/configs/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationArrayEndpoint(configurationService),
		rest.DecodeLoadConfigurationArrayRequest,
		common.EncodeResponse,
	))
}
