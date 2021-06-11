package main

import (
	"dictionaries-service/service"
	rest "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeConfigurationEndpoints(r *mux.Router, configurationService service.ConfigurationService) {

	// swagger:route GET /api/config/{key}/{day} loadConfig
	//
	// Loads all calendar items for given type and date range
	//     Responses:
	//       200: loadConfigurationOneResponseWrapper
	//       400: RestError
	r.Methods("GET").Path("/api/config/{key}/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationEndpoint(configurationService),
		rest.DecodeLoadConfigurationRequest,
		rest.EncodeResponse,
	))

	// swagger:route GET /api/configs/{day} loadConfigForKeys
	//
	// Loads all calendar items for given type and date range
	//     Responses:
	//       200: loadConfigurationResponse
	//       400: RestError
	r.Methods("GET").Path("/api/configs/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationArrayEndpoint(configurationService),
		rest.DecodeLoadConfigurationArrayRequest,
		rest.EncodeResponse,
	))
}
