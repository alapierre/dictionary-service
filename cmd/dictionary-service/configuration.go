package main

import (
	"dictionaries-service/service"
	rest "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeConfigurationEndpoints(r *mux.Router, configurationService service.ConfigurationService) {

	r.Methods("GET").Path("/api/config/{key}/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationEndpoint(configurationService),
		rest.DecodeLoadConfigurationRequest,
		rest.EncodeResponse,
	))

	r.Methods("GET").Path("/api/configs/{day}").Handler(http.NewServer(
		rest.MakeLoadConfigurationArrayEndpoint(configurationService),
		rest.DecodeLoadConfigurationArrayRequest,
		rest.EncodeResponse,
	))
}
