package main

import (
	"context"
	"dictionaries-service/service"
	rest "dictionaries-service/transport/http"
	"github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	nethttp "net/http"
)

func makeDictionariesEndpoints(r *mux.Router, dictionaryService *service.DictionaryService) {

	r.Methods("GET", "OPTIONS").Path("/api/dictionary/{type}/{key}").Handler(http.NewServer(
		rest.MakeLoadDictEndpoint(dictionaryService),
		rest.DecodeLoadDictRequest,
		rest.EncodeResponse,
	))

	r.Methods("GET", "OPTIONS").Path("/api/metadata/{type}").Handler(http.NewServer(
		rest.MakeLoadMetadataEndpoint(dictionaryService),
		rest.DecodeLoadMetadataRequest,
		rest.EncodeMetadataResponse,
	))

	r.Methods("GET", "OPTIONS").Path("/api/metadata").Handler(http.NewServer(
		rest.MakeAvailableDictionaryTypesEndpoint(dictionaryService),
		rest.EmptyRequest(),
		rest.EncodeResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/metadata").Handler(http.NewServer(
		rest.MakeSaveMetadataEndpoint(dictionaryService),
		rest.DecodeSaveMetadataRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/metadata/{type}").Handler(http.NewServer(
		rest.MakeSaveMetadataEndpointBetter(dictionaryService),
		rest.DecodeSaveMetadataRequestBetter,
		rest.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/metadata/{type}").Handler(http.NewServer(
		rest.MakeUpdateMetadataEndpointBetter(dictionaryService),
		rest.DecodeSaveMetadataRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/metadata/{type}").Handler(http.NewServer(
		rest.MakeSaveMetadataEndpointBetter(dictionaryService),
		rest.DecodeSaveMetadataRequestBetter,
		rest.EncodeSavedResponse,
	))

	r.Methods("GET", "OPTIONS").Path("/api/dictionary/{type}").Handler(http.NewServer(
		rest.MakeLoadDictionaryByType(dictionaryService),
		rest.DecodeByTypeRequest,
		rest.EncodeResponse,
	))

	r.Methods("GET", "OPTIONS").Path("/api/dictionary/{type}/{key}/shallow").Handler(http.NewServer(
		rest.MakeLoadDictShallowEndpoint(dictionaryService),
		rest.DecodeLoadDictRequest,
		rest.EncodeResponse,
	))

	r.Methods("GET", "OPTIONS").Path("/api/dictionary/{type}/{key}/children").Handler(http.NewServer(
		rest.MakeLoadDictChildrenEndpoint(dictionaryService),
		rest.DecodeLoadDictRequest,
		rest.EncodeResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/dictionary").Handler(http.NewServer(
		rest.MakeSaveDictionaryEndpoint(dictionaryService),
		rest.DecodeSaveDictRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/dictionary").Handler(http.NewServer(
		rest.MakeUpdateDictionaryEndpoint(dictionaryService),
		rest.DecodeSaveDictRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/dictionary/shallow").Handler(http.NewServer(
		rest.MakeShallowSaveDictionaryEndpoint(dictionaryService),
		rest.DecodeShallowSaveDictionaryRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/dictionary/shallow").Handler(http.NewServer(
		rest.MakeShallowUpdateDictionaryEndpoint(dictionaryService),
		rest.DecodeShallowSaveDictionaryRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("DELETE", "OPTIONS").Path("/api/dictionary/{type}/{key}").Handler(http.NewServer(
		rest.MakeDeleteDictionaryEndpoint(dictionaryService),
		rest.DecodeLoadDictRequest,
		rest.EncodeSavedResponse,
	))

	r.Methods("DELETE", "OPTIONS").Path("/api/dictionary/all").Handler(http.NewServer(
		rest.MakeDeleteAllDictionaryEndpoint(dictionaryService),
		func(ctx context.Context, request2 *nethttp.Request) (request interface{}, err error) {
			return nil, nil
		},
		rest.EncodeSavedResponse,
	))

	r.Methods("DELETE", "OPTIONS").Path("/api/dictionary/{type}").Handler(http.NewServer(
		rest.MakeDeleteDictionaryByTypeEndpoint(dictionaryService),
		rest.DecodeByTypeRequest,
		rest.EncodeSavedResponse,
	))
}
