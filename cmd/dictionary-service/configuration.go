package main

import (
	"dictionaries-service/configuration"
	rest "dictionaries-service/configuration/transport/http"
	common "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeConfigurationEndpoints(r *mux.Router, configurationService configuration.Service) {

	// swagger:route POST /api/config/value addNewConfigEntry
	//
	// Create new config value in time for existing key
	//     Responses:
	//       201:
	//       400: RestError
	r.Methods("POST", "OPTIONS").Path("/api/config/value").Handler(http.NewServer(
		rest.MakeAddNewConfigurationEntryEndpoint(configurationService),
		rest.DecodeAddNewConfigurationEntryRequest,
		common.EncodeWithStatus(201),
	))

	// swagger:route DELETE /api/config deleteConfigurationEntry
	//
	// Delete config value by key and date
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("DELETE", "OPTIONS").Path("/api/config").Handler(http.NewServer(
		rest.MakeDeleteConfigurationValueEndpoint(configurationService),
		rest.DecodeDeleteConfigurationRequest,
		common.EncodeWithStatus(204),
	))

	// swagger:route PUT /api/config updateConfiguration
	//
	// Save new config value
	//     Responses:
	//       204:
	//       400: RestError
	r.Methods("PUT", "OPTIONS").Path("/api/config").Handler(http.NewServer(
		rest.MakeUpdateConfigurationEndpoint(configurationService),
		rest.DecodeSaveConfigurationRequest,
		common.EncodeWithStatus(204),
	))

	// swagger:route POST /api/config saveConfiguration
	//
	// Save new config value
	//     Responses:
	//       201:
	//       400: RestError
	r.Methods("POST", "OPTIONS").Path("/api/config").Handler(http.NewServer(
		rest.MakeSaveConfigurationEndpoint(configurationService),
		rest.DecodeSaveConfigurationRequest,
		common.EncodeWithStatus(201),
	))

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
