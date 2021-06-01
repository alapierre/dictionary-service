package http

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"dictionaries-service/tenant"
	"encoding/json"
	"fmt"
	"github.com/go-eden/slf4go"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type saveMetadataRequest struct {
	Type    string
	Content string
}

type saveMetadataRequestBetter struct {
	Type    string                 `json:"type"`
	Content map[string]interface{} `json:"content"`
}

type metadataRequest struct {
	Type string
}

func MakeLoadMetadataEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(metadataRequest)
			res, err := service.LoadMetadata(req.Type, t.Name)
			if err != nil {
				return nil, err
			}
			return res.Content, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeSaveMetadataEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveMetadataRequest)
			err := service.SaveMetadata(metadataRequestToDictionaryMetadata(req, t.Name))
			if err != nil {
				return MakeRestError(err, "cant_create_new_dictionary_metadata")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeSaveMetadataEndpointBetter(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveMetadataRequestBetter)

			err := service.SaveMetadata(metadataRequestToDictionaryMetadataBetter(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_create_new_dictionary_metadata")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func MakeUpdateMetadataEndpointBetter(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if t, ok := tenant.FromContext(ctx); ok {
			req := request.(saveMetadataRequestBetter)
			slog.Info("Trying to save: req: ", request)

			err := service.UpdateMetadata(metadataRequestToDictionaryMetadataBetter(req, t.Name))

			if err != nil {
				return MakeRestError(err, "cant_update_dictionary_metadata")
			}
			return nil, nil
		}
		return MakeRestError(fmt.Errorf("can't extract tenant from context"), "cant_extract_tenant_from_context")
	}
}

func DecodeLoadMetadataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dictionaryType := vars["type"]
	return metadataRequest{Type: dictionaryType}, nil
}

func DecodeSaveMetadataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var request saveMetadataRequest
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeSaveMetadataRequestBetter(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	dictionaryType := vars["type"]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var request saveMetadataRequestBetter
	request.Type = dictionaryType

	if err := json.Unmarshal(bodyBytes, &request.Content); err != nil {
		return nil, err
	}

	return request, nil
}

func EncodeMetadataResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
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

	_, err := io.WriteString(w, response.(string))
	return err
}

func metadataRequestToDictionaryMetadata(req saveMetadataRequest, tenant string) *model.DictionaryMetadata {
	return &model.DictionaryMetadata{
		Type:    req.Type,
		Tenant:  tenant,
		Content: req.Content,
	}
}

func metadataRequestToDictionaryMetadataBetter(req saveMetadataRequestBetter, tenant string) *model.DictionaryMetadata {

	content, _ := json.Marshal(req.Content)

	return &model.DictionaryMetadata{
		Type:    req.Type,
		Tenant:  tenant,
		Content: string(content),
	}
}
