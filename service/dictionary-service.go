package service

import (
	"dictionaries-service/model"
	slog "github.com/go-eden/slf4go"
	"golang.org/x/text/language"
)

func NewDictionaryService(dictionaryRepository DictionaryRepository, translateRepository TranslateRepository, metadataRepository DictionaryMetadataRepository) *DictionaryService {
	return &DictionaryService{
		dictionaryRepository:  dictionaryRepository,
		translationRepository: translateRepository,
		metadataRepository:    metadataRepository,
	}
}

type DictionaryService struct {
	dictionaryRepository  DictionaryRepository
	translationRepository TranslateRepository
	metadataRepository    DictionaryMetadataRepository
}

func (s *DictionaryService) LoadShallow(key, dictionaryType, tenant string) (map[string]interface{}, error) {
	dict, err := s.dictionaryRepository.Load(key, dictionaryType, tenant)
	if err != nil {
		return nil, err
	}
	return prepareMap(dict), nil
}

func (s *DictionaryService) LoadShallowTranslated(key, dictionaryType, tenant string, lang language.Tag) (map[string]interface{}, error) {
	dict, err := s.dictionaryRepository.LoadTranslated(key, dictionaryType, tenant, lang)
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

func (s *DictionaryService) LoadTranslated(key, dictionaryType, tenant string, lang language.Tag) (map[string]interface{}, error) {
	dict, err := s.dictionaryRepository.LoadTranslated(key, dictionaryType, tenant, lang)
	if err != nil {
		slog.Error("can't load translated parent dictionary ", err)
		return nil, err
	}

	res := prepareMap(dict)

	if dict.ParentKey == nil {
		children, err := s.dictionaryRepository.LoadChildrenDictionaryTranslated(key, dictionaryType, tenant, lang)
		if err != nil {
			return nil, err
		}
		res["children"] = prepareDictionaryChildrenMap(children, false)
	}

	return res, nil
}

func (s *DictionaryService) LoadChildrenTranslated(key, dictionaryType, tenant string, lang language.Tag) ([]childrenMap, error) {

	children, err := s.dictionaryRepository.LoadChildrenDictionaryTranslated(key, dictionaryType, tenant, lang)
	if err != nil {
		return nil, err
	}

	res := prepareDictionaryChildrenMap(children, true)

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

// jeśli parametrem jest parent, zapisuje jeden poziom słownika. bez wpływu na dzieci
func (s *DictionaryService) SaveShallow(dict *model.Dictionary) error {
	return s.dictionaryRepository.Save(dict)
}

// jeśli parametrem jest parent, aktualizuje jeden poziom słownika. bez wpływu na dzieci
func (s *DictionaryService) UpdateShallow(dict *model.Dictionary) error {
	return s.dictionaryRepository.Update(dict)
}

func (s *DictionaryService) Delete(key, dictionaryType, tenant string) error {
	return s.dictionaryRepository.Delete(key, dictionaryType, tenant)
}

func (s *DictionaryService) DeleteAll(tenant string) error {
	return s.dictionaryRepository.DeleteAll(tenant)
}

func (s *DictionaryService) DeleteByType(dictionaryType, tenant string) error {
	return s.dictionaryRepository.DeleteByType(dictionaryType, tenant)
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
		GroupId:   parent.GroupId,
	}

	return s.dictionaryRepository.Save(&dict)
}

func (s *DictionaryService) AvailableTranslation(key, dictionaryType, tenant string) ([]string, error) {
	return s.translationRepository.AvailableTranslation(key, dictionaryType, tenant)
}

func (s *DictionaryService) LoadMetadata(dictionaryType, tenant string) (*model.DictionaryMetadata, error) {
	return s.metadataRepository.Load(dictionaryType, tenant)
}

func (s *DictionaryService) SaveMetadata(metadata *model.DictionaryMetadata) error {
	return s.metadataRepository.Save(metadata)
}

func (s *DictionaryService) UpdateMetadata(metadata *model.DictionaryMetadata) error {
	return s.metadataRepository.Update(metadata)
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
		return s.dictionaryRepository.DeleteMultiple(keys, parent.Tenant, parent.Type)
	}
	return nil
}

func (s *DictionaryService) LoadByType(dictionaryType string, tenant string) ([]model.Dictionary, error) {
	return s.dictionaryRepository.LoadByType(dictionaryType, tenant)
}

func (s *DictionaryService) AvailableDictionaryTypes(tenant string) ([]string, error) {
	return s.metadataRepository.AvailableDictionaryTypes(tenant)
}

type childrenMap map[string]interface{}

func prepareMap(dict *model.Dictionary) map[string]interface{} {

	res := make(map[string]interface{})

	res["key"] = dict.Key
	res["type"] = dict.Type
	res["name"] = dict.Name

	if dict.GroupId != nil {
		res["group_id"] = dict.GroupId
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

		mergeMaps(&ch.Content, &tmp)
		list[i] = tmp
	}
	return list
}

func prepareDictionaryChildrenMap(children []model.ChildDictionary, withParentKey bool) []childrenMap {

	list := make([]childrenMap, len(children))

	for i, ch := range children {
		tmp := make(map[string]interface{})
		tmp["key"] = ch.Key
		tmp["name"] = ch.Name
		if withParentKey {
			tmp["parent_key"] = ch.ParentKey
		}
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
		GroupId: parent.GroupId,
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
		GroupId:   parent.GroupId,
		Tenant:    parent.Tenant,
		Content:   child.Content,
		ParentKey: &parent.Key,
	}
	return toSave
}
