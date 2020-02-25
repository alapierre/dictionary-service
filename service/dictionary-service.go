package service

import (
	"dictionaries-service/model"
)

func NewDictionaryService(dictionaryRepository DictionaryRepository) *DictionaryService {
	return &DictionaryService{dictionaryRepository: dictionaryRepository}
}

type DictionaryService struct {
	dictionaryRepository DictionaryRepository
}

func (s *DictionaryService) LoadShallow(key, dictionaryType, tenant string) (map[string]interface{}, error) {
	dict, err := s.dictionaryRepository.Load(key, dictionaryType, tenant)
	if err != nil {
		return nil, err
	}
	return prepareMap(dict), nil
}

func (s *DictionaryService) Load(key, dictionaryType, tenant string) (map[string]interface{}, error) {
	dict, err := s.dictionaryRepository.Load(key, dictionaryType, tenant)
	if err != nil {
		return nil, err
	}

	res := prepareMap(dict)

	if dict.ParentKey == nil {
		children, err := s.dictionaryRepository.LoadChildren(key, dictionaryType, tenant)
		if err != nil {
			return nil, err
		}
		res["children"] = prepareChildrenMap(children)
	}

	return res, nil
}

type childrenMap map[string]interface{}

func prepareMap(dict *model.Dictionary) map[string]interface{} {

	res := make(map[string]interface{})

	res["key"] = dict.Key
	res["type"] = dict.Type
	res["tenant"] = dict.Tenant
	res["name"] = dict.Name

	if dict.GroupID != nil {
		res["group_id"] = dict.GroupID
	}

	if dict.ParentKey != nil {
		res["parent_key"] = *dict.ParentKey
	}

	mergeMaps(&dict.Content, &res)

	return res
}

func prepareChildrenMap(children []model.Dictionary) []childrenMap {

	list := make([]childrenMap, len(children))

	for i, ch := range children {
		tmp := make(map[string]interface{})
		tmp["type"] = ch.Type
		tmp["key"] = ch.Key
		tmp["name"] = ch.Name
		//tmp["tenant"] = ch.Tenant

		mergeMaps(&ch.Content, &tmp)
		list[i] = tmp
	}
	return list
}

func mergeMaps(source, destination *map[string]interface{}) {
	for k, v := range *source {
		(*destination)[k] = v
	}
}
