package http

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"dictionaries-service/tenant"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"io/ioutil"
	"net/http"
)

var DefaultLanguage language.Tag

// RestError general service error message
// swagger:response RestError
type RestError struct {

	// Error message
	Error string `json:"error,omitempty"`
	// An optional field describes error in more details
	ErrorDescription string `json:"error_description,omitempty"`
}

type dictionaryRequest struct {
	Key  string
	Type string
}

type saveDictionaryRequest struct {
	Key      string
	Type     string
	Name     string
	GroupId  *string
	Content  map[string]interface{}
	Children []map[string]interface{}
}

type byTypeRequest struct {
	Type string
}

type saveShallowDictionaryRequest struct {
	Key       string                 `json:"key"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	GroupId   *string                `json:"group_id"`
	Content   map[string]interface{} `json:"content"`
	ParentKey *string                `json:"parent_key"`
}

func MakeLoadDictShallowEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		if t, ok := tenant.FromContext(ctx); ok {

			req := request.(dictionaryRequest)
			r, err := service.LoadShallow(req.Key, req.Type, t.Name)

			if err != nil {
				return MakeRestError(err, "cant_load_dictionary_by_key_and_type")
			}
			return r, nil

		} else {
			return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}
	}
}

func MakeLoadDictChildrenEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(dictionaryRequest)
			lang := extractLang(service, ctx, req.Key, req.Type, t.Name)

			r, err := service.LoadChildrenTranslated(req.Key, req.Type, t.Name, lang)

			if err != nil {
				return MakeRestError(err, "cant_load_dictionary_by_key_and_type")
			}
			return r, nil
		} else {
			return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
		}
	}
}

func MakeLoadDictionaryByType(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(byTypeRequest)
			res, err := service.LoadByType(req.Type, t.Name)
			return res, err
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeAvailableDictionaryTypesEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			types, err := service.AvailableDictionaryTypes(t.Name)
			return types, err
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeDeleteDictionaryByTypeEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(byTypeRequest)
			return nil, service.DeleteByType(req.Type, t.Name)
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func DecodeByTypeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dictionaryType := vars["type"]
	return byTypeRequest{Type: dictionaryType}, nil
}

func MakeDeleteAllDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			return nil, service.DeleteAll(t.Name)
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeDeleteDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(dictionaryRequest)
			return nil, service.Delete(req.Key, req.Type, t.Name)
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeShallowUpdateDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveShallowDictionaryRequest)

			err := service.UpdateShallow(shallowDictionaryToDictionary(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_create_new_dictionary_entry")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeShallowSaveDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveShallowDictionaryRequest)

			err := service.SaveShallow(shallowDictionaryToDictionary(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_create_new_dictionary_entry")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func DecodeShallowSaveDictionaryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request saveShallowDictionaryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func MakeSaveDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveDictionaryRequest)
			err := service.SaveParent(convertRequestToDictionary(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_create_new_dictionary_entry")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeUpdateDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveDictionaryRequest)
			err := service.UpdateParent(convertRequestToDictionary(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_update_dictionary_entry_by_key_and_type")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeLoadDictEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(dictionaryRequest)
			lang := extractLang(service, ctx, req.Key, req.Type, t.Name)
			r, err := service.LoadTranslated(req.Key, req.Type, t.Name, lang)

			if err != nil {
				return MakeRestError(err, "cant_load_dictionary_by_key_and_type")
			}
			return r, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func DecodeSaveDictRequest(_ context.Context, r *http.Request) (interface{}, error) {

	// TODO: go away with double json parse

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var request saveDictionaryRequest
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		return nil, err
	}

	var content map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &content); err != nil {
		return nil, err
	}

	delete(content, "key")      // because we not need dictionary key in content
	delete(content, "name")     // because we not need dictionary name in content
	delete(content, "type")     // because we not need dictionary type in content
	delete(content, "tenant")   // because we not need tenant in dictionary content
	delete(content, "children") // because we not need children in dictionary content

	request.Content = content
	return request, nil
}

func DecodeLoadDictRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	key := vars["key"]
	dictionaryType := vars["type"]
	return dictionaryRequest{Key: key, Type: dictionaryType}, nil
}

func shallowDictionaryToDictionary(req saveShallowDictionaryRequest, tenant string) *model.Dictionary {
	return &model.Dictionary{
		Key:       req.Key,
		Type:      req.Type,
		Name:      req.Name,
		GroupId:   req.GroupId,
		Tenant:    tenant,
		Content:   req.Content,
		ParentKey: req.ParentKey,
	}
}
