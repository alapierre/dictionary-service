package http

import (
	"context"
	"dictionaries-service/service"
	"dictionaries-service/tenant"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type configurationArrayRequest struct {
	Keys []string
	Day  time.Time
}

func MakeLoadConfigurationArrayEndpoint(configurationService service.ConfigurationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
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

type configurationRequest struct {
	Key string
	Day time.Time
}

type loadConfigurationResponse struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
	Type  string  `json:"type"`
}

func MakeLoadConfigurationEndpoint(configurationService service.ConfigurationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}

		req := request.(configurationRequest)

		r, err := configurationService.LoadForDay(req.Key, t.Name, req.Day)

		if err != nil {
			return MakeRestError(err, "cant_load_configuration_by_key_tenant_and_day")
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
