package transport

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"encoding/json"
	slog "github.com/go-eden/slf4go"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

var DefaultLanguage language.Tag

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

type byTypeRequest struct {
	Type string
}

type metadataRequest struct {
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

		tenant := extractTenant(ctx)
		req := request.(dictionaryRequest)

		r, err := service.LoadShallow(req.Key, req.Type, tenant)

		if err != nil {
			return makeRestError(err, "cant_load_dictionary_by_key_and_type")
		}
		return r, nil
	}
}

func MakeLoadDictionaryByType(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(byTypeRequest)
		res, err := service.LoadByType(req.Type, tenant)
		return res, err
	}
}

func MakeLoadMetadataEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(metadataRequest)
		res, err := service.LoadMetadata(req.Type, tenant)
		if err != nil {
			return nil, err
		}
		return res.Content, nil
	}
}

func MakeAvailableDictionaryTypesEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		types, err := service.AvailableDictionaryTypes(tenant)
		return types, err
	}
}

func DecodeLoadMetadataRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dictionaryType := vars["type"]
	return metadataRequest{Type: dictionaryType}, nil
}

//func MakeTranslateDictionaryItemEndpoint(service *service.DictionaryService) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (interface{}, error) {
//
//	}
//}

func MakeDeleteDictionaryByTypeEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(byTypeRequest)
		return nil, service.DeleteByType(req.Type, tenant)
	}
}

func DecodeByTypeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	dictionaryType := vars["type"]
	return byTypeRequest{Type: dictionaryType}, nil
}

func MakeDeleteAllDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		return nil, service.DeleteAll(tenant)
	}
}

func MakeDeleteDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(dictionaryRequest)
		return nil, service.Delete(req.Key, req.Type, tenant)
	}
}

func MakeShallowUpdateDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(saveShallowDictionaryRequest)

		err := service.UpdateShallow(shallowDictionaryToDictionary(req, tenant))

		if err != nil {
			return makeRestError(err, "cant_create_new_dictionary_entry")
		}
		return nil, nil
	}
}

func MakeShallowSaveDictionaryEndpoint(service *service.DictionaryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tenant := extractTenant(ctx)
		req := request.(saveShallowDictionaryRequest)

		err := service.SaveShallow(shallowDictionaryToDictionary(req, tenant))

		if err != nil {
			return makeRestError(err, "cant_create_new_dictionary_entry")
		}
		return nil, nil
	}
}

func DecodeShallowSaveDictionaryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request saveShallowDictionaryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
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
		lang := extractLang(service, ctx, req.Key, req.Type, tenant)

		r, err := service.LoadTranslated(req.Key, req.Type, tenant, lang)

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

func extractLang(translator service.Translator, ctx context.Context, key, dictionaryType, tenant string) language.Tag {

	var matcher = language.NewMatcher(loadAvailableTransitions(translator, key, dictionaryType, tenant))

	accept := ctx.Value("language").(string)
	tag, _ := language.MatchStrings(matcher, accept)

	return tag
}

func loadAvailableTransitions(translator service.Translator, key, dictionaryType, tenant string) []language.Tag {

	var tags []language.Tag
	tags = append(tags, DefaultLanguage) // The first language is used as fallback.

	base, _ := DefaultLanguage.Base()
	languageToSkip := base.String()

	tr, err := translator.AvailableTranslation(key, dictionaryType, tenant)

	if err != nil {
		slog.Warn("Can't load translations")
	}

	for _, i := range tr {
		if i != languageToSkip { // Default already added
			tag, err := language.Parse(i)
			if err != nil {
				tags = append(tags, tag)
			}
		}
	}

	return tags
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
