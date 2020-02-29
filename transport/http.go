package transport

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"io/ioutil"
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

type saveDictionaryRequest struct {
	Key      string
	Type     string
	Name     string
	GroupId  *string
	Content  map[string]interface{}
	Children []map[string]interface{}
}

type saveDictionaryChildRequest struct {
	ParentKey string
	Key       string
	Name      string
	Content   map[string]interface{}
}

func MakeSaveDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(saveDictionaryRequest)

		err := service.SaveParent(convertRequestToDictionary(req, tenant))

		if err != nil {
			return makeRestError(err, "cant_create_new_dictionary_entry")
		}
		return nil, nil
	}
}

func MakeUpdateDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(saveDictionaryRequest)

		err := service.UpdateParent(convertRequestToDictionary(req, tenant))

		if err != nil {
			return makeRestError(err, "cant_update_dictionary_entry_by_key_and_type")
		}
		return nil, nil
	}
}

func convertRequestToDictionary(req saveDictionaryRequest, tenant string) *model.ParentDictionary {
	var children []model.ChildDictionary

	for _, c := range req.Children {
		children = append(children, model.ChildDictionary{
			Key:       c["key"].(string),
			Name:      c["name"].(string),
			ParentKey: req.Key,
			Content:   c,
		})
		delete(c, "name") // because we not need child name in content
		delete(c, "key")  // because we not need child name in content
		delete(c, "type") // because we not need child name in content
	}

	s := model.ParentDictionary{
		Key:      req.Key,
		Type:     req.Type,
		Name:     req.Name,
		GroupId:  req.GroupId,
		Tenant:   tenant,
		Content:  req.Content,
		Children: children,
	}
	return &s
}

func DecodeSaveDictRequest(_ context.Context, r *http.Request) (interface{}, error) {

	// TODO: do away with double json parse

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
	delete(content, "type")     // because we not need dictionary name in content
	delete(content, "tenant")   // because we not need tenant in dictionary content
	delete(content, "children") // because we not need children in dictionary content

	request.Content = content
	return request, nil
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

func DecodeLoadDictRequest(_ context.Context, r *http.Request) (interface{}, error) {
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

func EncodeSavedResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
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
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	return json.NewEncoder(w).Encode(response)
}
