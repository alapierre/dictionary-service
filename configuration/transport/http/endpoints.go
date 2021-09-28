package http

import (
	"context"
	"dictionaries-service/configuration"
	"dictionaries-service/tenant"
	commons "dictionaries-service/transport/http"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// swagger:parameters loadConfigForKeys
type configurationArrayRequest struct {
	Keys []string `json:"key"`

	// swagger:strfmt date
	// in:path
	Day time.Time `json:"day"`
}

// swagger:parameters loadConfig
type configurationRequest struct {
	// in:path
	Key string `json:"key"`
	// in:path
	Day time.Time `json:"day"`
}

type loadConfigurationResponse struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
	Type  string  `json:"type"`
}

// swagger:response loadConfigurationResponse
//goland:noinspection ALL
type loadConfigurationResponseWrapper struct {

	// in:body
	Body []loadConfigurationResponse
}

// swagger:response loadShortResponseWrapper
//goland:noinspection ALL
type loadShortResponseWrapper struct {

	// in:body
	Body []configuration.Short
}

// swagger:response loadConfigurationOneResponseWrapper
//goland:noinspection ALL
type loadConfigurationOneResponseWrapper struct {

	// in:body
	Body loadConfigurationResponse
}

func MakeLoadAllShortEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		short, err := configurationService.LoadAllShort(ctx)
		if err != nil {
			return commons.MakeRestError(err, "cant_delete_dictionary_entry")
		}

		return short, nil
	}
}

func MakeLoadConfigurationArrayEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return commons.MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}

		req := request.(configurationArrayRequest)

		configs := configurationService.LoadMany(t.Name, req.Day, req.Keys...)
		var res []loadConfigurationResponse

		var value *string

		for _, t := range configs {

			if t.Value.Valid {
				tmp := t.Value.String // I need new copy of returned string
				value = &tmp          // and than pointer to it
			} else {
				continue // skip null values in result
			}

			res = append(res, loadConfigurationResponse{
				Key:   t.Key,
				Value: value,
				Type:  t.Type,
			})
		}

		return res, nil
	}
}

func DecodeLoadConfigurationArrayRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	day, err := time.Parse("2006-01-02", vars["day"])

	if err != nil {
		return nil, err
	}

	return configurationArrayRequest{
		Keys: r.URL.Query()["key"],
		Day:  day,
	}, nil
}

func MakeLoadConfigurationEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return commons.MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}

		req := request.(configurationRequest)

		r, err := configurationService.LoadForDay(req.Key, t.Name, req.Day)

		if err != nil {
			return commons.MakeRestError(err, "cant_load_configuration_by_key_tenant_and_day")
		}

		var value *string

		if r.Value.Valid {
			value = &r.Value.String
		} else {
			value = nil
		}

		return &loadConfigurationResponse{
			Key:   r.Key,
			Value: value,
			Type:  r.Type,
		}, nil
	}
}

func DecodeLoadConfigurationRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	day, err := time.Parse("2006-01-02", vars["day"])

	if err != nil {
		return nil, err
	}

	return configurationRequest{
		Key: vars["key"],
		Day: day,
	}, nil
}
