package http

import (
	"context"
	"dictionaries-service/configuration"
	"dictionaries-service/tenant"
	commons "dictionaries-service/transport/http"
	"dictionaries-service/types"
	"dictionaries-service/util"
	"encoding/json"
	"fmt"
	slog "github.com/go-eden/slf4go"
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

// swagger:parameters loadValues
type valuesRequest struct {

	// in:path
	Key string `json:"key"`

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`
}

// swagger:parameters loadConfig
type configurationRequest struct {
	// in:path
	Key string `json:"key"`
	// in:path
	Day time.Time `json:"day"`
}

type saveConfigurationRequest struct {
	Key      string         `json:"key"`
	Value    *string        `json:"value"`
	Type     string         `json:"type"`
	Name     string         `json:"name"`
	DateFrom types.JsonDate `json:"date_from"`
	DateTo   types.JsonDate `json:"date_to"`
}

type addNewConfigurationValueRequest struct {
	Value    *string        `json:"value"`
	DateFrom types.JsonDate `json:"date_from"`
	DateTo   types.JsonDate `json:"date_to"`
}

type addNewConfigurationValueCompleted struct {
	Key string `json:"key"`
	addNewConfigurationValueRequest
}

type deleteConfigurationRequest struct {
	Key      string         `json:"key"`
	DateFrom types.JsonDate `json:"date_from"`
}

type loadValueResponseFull struct {
	Key    string        `json:"key"`
	Type   string        `json:"type"`
	Name   string        `json:"name"`
	Values []valueInTime `json:"values"`
}

type valueInTime struct {
	Value    *string        `json:"value"`
	DateFrom types.JsonDate `json:"date_from"`
	DateTo   types.JsonDate `json:"date_to"`
}

type loadConfigurationResponse struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
	Type  string  `json:"type"`
}

type updateValueInTimeRequest struct {
	Value *string `json:"value"`
}

type updateValueInTimeCompleted struct {
	key  string
	from time.Time
	updateValueInTimeRequest
}

func DecodeUpdateValueInTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var request updateValueInTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	key := vars["key"]
	from, err := util.StringToDate(vars["from"])

	if err != nil {
		slog.Errorf("can't convert %s into date", vars["from"])
		return nil, err
	}

	return updateValueInTimeCompleted{
		key:                      key,
		from:                     from,
		updateValueInTimeRequest: request,
	}, nil
}

func MakeUpdateValueInTimeEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(updateValueInTimeCompleted)
		err := configurationService.UpdateValueInTime(ctx, req.key, req.from, req.Value)
		if err != nil {
			return commons.MakeRestError(fmt.Errorf("can't update in time value for given key: %s, and date: %s, %v", req.key, req.from, err),
				"cant_add_new_in_time_entry")
		}
		return nil, nil
	}
}

func DecodeAddNewConfigurationValueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request addNewConfigurationValueRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	vars := mux.Vars(r)
	key := vars["key"]

	return addNewConfigurationValueCompleted{
		Key:                             key,
		addNewConfigurationValueRequest: request,
	}, nil
}

func MakeAddNewConfigurationValueEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(addNewConfigurationValueCompleted)
		if err := configurationService.AddNewValueInTime(ctx, req.Key, req.Value, req.DateFrom.Time(), req.DateTo.Time()); err != nil {
			return commons.MakeRestError(fmt.Errorf("can't create new in time value for given key: %s, and date: %s, %v", req.Key, req.DateFrom, err),
				"cant_add_new_in_time_value")
		}
		return nil, nil
	}
}

func DecodeDeleteConfigurationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteConfigurationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func MakeDeleteConfigurationValueEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(deleteConfigurationRequest)
		if err := configurationService.DeleteValue(ctx, req.Key, req.DateFrom.Time()); err != nil {
			return commons.MakeRestError(fmt.Errorf("can't delete entry for given key: %s, and date: %s", req.Key, req.DateFrom),
				"cant_delete")
		}
		return nil, nil
	}
}

func MakeUpdateConfigurationEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(saveConfigurationRequest)

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return commons.MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}

		err := configurationService.Update(&configuration.Configuration{
			Key:      req.Key,
			Tenant:   t.Name,
			Type:     req.Type,
			Name:     req.Name,
			Value:    util.PointerToSqlNullString(req.Value),
			DateFrom: req.DateFrom.Time(),
			DateTo:   req.DateTo.Time(),
		})

		if err != nil {
			return commons.MakeRestError(err, "cant_update_config_entry")
		}

		return nil, err
	}
}

func MakeSaveConfigurationEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(saveConfigurationRequest)

		t, ok := tenant.FromContext(ctx)
		if !ok {
			return commons.MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}

		err := configurationService.Save(&configuration.Configuration{
			Key:      req.Key,
			Tenant:   t.Name,
			Type:     req.Type,
			Name:     req.Name,
			Value:    util.PointerToSqlNullString(req.Value),
			DateFrom: req.DateFrom.Time(),
			DateTo:   req.DateTo.Time(),
		})

		if err != nil {
			return commons.MakeRestError(err, "cant_save_config_entry")
		}

		return nil, err
	}
}

func DecodeSaveConfigurationRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var request saveConfigurationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
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

func MakeLoadValuesEndpoint(configurationService configuration.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(valuesRequest)

		entry, err := configurationService.LoadEntry(ctx, req.Key)
		if err != nil {
			return commons.MakeRestError(err, "cant_delete_dictionary_entry")
		}

		values, err := configurationService.LoadValues(ctx, req.Key)
		if err != nil {
			return commons.MakeRestError(err, "cant_delete_dictionary_entry")
		}

		var res []valueInTime

		for _, v := range values {
			res = append(res, valueInTime{
				Value:    util.SqlNullStringToStringPointer(v.Value),
				DateFrom: types.JsonDate(v.DateFrom),
				DateTo:   types.JsonDate(v.DateTo),
			})
		}

		return &loadValueResponseFull{
			Key:    entry.Key,
			Type:   entry.Type,
			Name:   entry.Name,
			Values: res,
		}, nil
	}
}

func DecodeLoadValuesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return valuesRequest{Key: vars["key"]}, nil
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
