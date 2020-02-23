package transport

import (
	"context"
	"dictionaries-service/service"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type RestError struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type dictionaryRequest struct {
	Key  string
	Type string
}

func MakeLoadDictEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		tenant := extractTenant(ctx)
		req := request.(dictionaryRequest)
		r, err := service.Load(req.Key, req.Type, tenant)

		if err != nil {
			return makeRestError(err, "cant_load_dictionary_by_key_and_type")
		}
		return r, nil
	}
}

func LoadDictRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key := vars["key"]
	dictionaryType := vars["type"]
	return dictionaryRequest{Key: key, Type: dictionaryType}, nil
}

func extractTenant(ctx context.Context) string {
	return ctx.Value("tenant").(string)
}

func makeRestError(err error, message string) (interface{}, error) {
	return &RestError{
		Error:            message,
		ErrorDescription: err.Error(),
	}, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	headers := w.Header()
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	headers.Set("X-Content-Type-Options", "nosniff")
	headers.Set("X-XSS-Protection", "1; mode=block")
	headers.Set("Pragma", "no-cache")
	headers.Set("Expires", "0")
	headers.Set("X-Frame-Options", "DENY")

	if _, err := response.(*RestError); err {
		w.WriteHeader(http.StatusBadRequest)
	}

	if reflect.ValueOf(response).IsNil() {
		rt := reflect.TypeOf(response)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			return json.NewEncoder(w).Encode(make([]int, 0))
		}
	}

	return json.NewEncoder(w).Encode(response)
}
