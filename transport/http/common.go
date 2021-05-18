package http

import (
	"context"
	"dictionaries-service/model"
	"dictionaries-service/service"
	"github.com/go-eden/slf4go"
	"golang.org/x/text/language"
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

func makeRestError(err error, message string) (interface{}, error) {
	return &RestError{
		Error:            message,
		ErrorDescription: err.Error(),
	}, nil
}
