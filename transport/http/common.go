package http

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"encoding/json"
	"github.com/go-eden/slf4go"
	"golang.org/x/text/language"
	"net/http"
	"reflect"
)

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

func MakeRestError(err error, message string) (interface{}, error) {
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
