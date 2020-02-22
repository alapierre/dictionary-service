package service

import (
	"dictionaries-service/model"
)

func NewDictionaryService(dictionaryRepository DictionaryRepository, childRepository ChildRepository) *DictionaryService {
	return &DictionaryService{dictionaryRepository: dictionaryRepository, childRepository: childRepository}
}

type DictionaryService struct {
	dictionaryRepository DictionaryRepository
	childRepository      ChildRepository
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

	if dict.Parent {
		children, err := s.childRepository.LoadChildren(key, dictionaryType, tenant)
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
	res["group_id"] = dict.GroupID
	res["parent_key"] = dict.ParentKey
	mergeMaps(&dict.Content, &res)

	return res
}

func prepareChildrenMap(children []model.Child) []childrenMap {

	list := make([]childrenMap, len(children))

	for i, ch := range children {
		tmp := make(map[string]interface{})
		tmp["type"] = ch.Type
		tmp["key"] = ch.Key
		tmp["tenant"] = ch.Tenant

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
