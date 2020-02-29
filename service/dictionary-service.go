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

func (s *DictionaryService) SaveParent(parent *model.ParentDictionary) error {
	return s.dictionaryRepository.DoInTransaction(func() error {

		dict := parentToDictionary(parent)

		err := s.dictionaryRepository.Save(&dict)
		if err != nil {
			return err
		}

		if len(parent.Children) != 0 {
			return s.saveChildren(parent)
		}

		return nil
	})
}

func (s *DictionaryService) UpdateParent(parent *model.ParentDictionary) error {
	return s.dictionaryRepository.DoInTransaction(func() error {

		dict := parentToDictionary(parent)
		err := s.dictionaryRepository.Update(&dict)
		if err != nil {
			return err
		}

		if len(parent.Children) != 0 {
			return s.updateChildren(parent)
		}

		return nil
	})
}

func (s *DictionaryService) SaveChild(child *model.ChildDictionary, dictionaryType, tenant string) error {

	parent, err := s.dictionaryRepository.Load(child.ParentKey, dictionaryType, tenant)
	if err != nil {
		return err
	}

	dict := model.Dictionary{
		Key:       child.Key,
		Name:      child.Name,
		Content:   child.Content,
		ParentKey: &child.ParentKey,
		GroupID:   parent.GroupID,
	}

	return s.dictionaryRepository.Save(&dict)
}

// transaction should be began before call this
func (s *DictionaryService) saveChildren(parent *model.ParentDictionary) error {
	for _, child := range parent.Children {
		toSave := childToDictionary(child, parent)
		err := s.dictionaryRepository.Save(&toSave)
		if err != nil {
			return err
		}
	}
	return nil
}

// transaction should be began before call this
func (s *DictionaryService) updateChildren(parent *model.ParentDictionary) error {

	keys, err := s.dictionaryRepository.LoadChildrenKeys(parent.Key, parent.Type, parent.Tenant)

	for _, child := range parent.Children {
		toSave := childToDictionary(child, parent)
		if keys[child.Key] {
			err = s.dictionaryRepository.Update(&toSave)
			delete(keys, child.Key)
		} else {
			err = s.dictionaryRepository.Save(&toSave)
		}
		if err != nil {
			return err
		}
	}
	if len(keys) != 0 {
		return s.dictionaryRepository.DeleteMultiple(keys)
	}
	return nil
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

func parentToDictionary(parent *model.ParentDictionary) model.Dictionary {
	dict := model.Dictionary{
		Key:     parent.Key,
		Type:    parent.Type,
		Name:    parent.Name,
		GroupID: parent.GroupId,
		Tenant:  parent.Tenant,
		Content: parent.Content,
	}
	return dict
}

func childToDictionary(child model.ChildDictionary, parent *model.ParentDictionary) model.Dictionary {
	toSave := model.Dictionary{
		Key:       child.Key,
		Type:      parent.Type,
		Name:      child.Name,
		GroupID:   parent.GroupId,
		Tenant:    parent.Tenant,
		Content:   child.Content,
		ParentKey: &parent.Key,
	}
	return toSave
}
