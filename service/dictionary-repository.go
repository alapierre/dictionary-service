package service

import (
	"dictionaries-service/model"
	"dictionaries-service/util"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v10"
	"golang.org/x/text/language"
)

func NewDictionaryRepository(db *pg.DB) DictionaryRepository {
	return &dictionaryRepository{db: db}
}

type DictionaryRepository interface {
	Load(key, dictionaryType, tenant string) (*model.Dictionary, error)
	Save(dict *model.Dictionary) error
	Update(dict *model.Dictionary) error
	LoadAll(tenant string) ([]model.Dictionary, error)
	LoadByType(dictionaryType, tenant string) ([]model.Dictionary, error)
	LoadChildren(parentKey, dictionaryType, tenant string) ([]model.Dictionary, error)
	DoInTransaction(callback util.TransactionCallback) error
	LoadChildrenKeys(parentKey, dictionaryType, tenant string) (ChildrenKeys, error)
	DeleteMultiple(keys ChildrenKeys, tenant, dictionaryType string) error
	Delete(key, dictionaryType, tenant string) error
	DeleteByType(dictionaryType, tenant string) error
	DeleteAll(tenant string) error
	LoadTranslated(key, dictionaryType, tenant string, lang language.Tag) (*model.Dictionary, error)
	LoadChildrenTranslated(parentKey, dictionaryType, tenant string, lang language.Tag) ([]model.Dictionary, error)
	LoadChildrenDictionaryTranslated(parentKey, dictionaryType, tenant string, lang language.Tag) ([]model.ChildDictionary, error)
}

type dictionaryRepository struct {
	db *pg.DB
}

func (s *dictionaryRepository) DoInTransaction(callback util.TransactionCallback) error {
	return util.DoInTransaction(s.db, callback)
}

func (s *dictionaryRepository) Load(key, dictionaryType, tenant string) (*model.Dictionary, error) {

	var dict model.Dictionary

	_, err := s.db.QueryOne(&dict,
		`select * from all_dictionaries where tenant = ? and key = ? and type = ?`, tenant, key, dictionaryType)

	return &dict, err
}

func (s *dictionaryRepository) LoadTranslated(key, dictionaryType, tenant string, lang language.Tag) (*model.Dictionary, error) {

	var dict model.Dictionary

	sql := `select d.key, d.type, COALESCE(t.name, d.name) as name, group_id, d.tenant, content, parent, parent_key 
				from all_dictionaries d
    			  left outer join translation t on d.key = t.key and d.tenant = t.tenant and d.type = t.type
                where (t.language = ? or t.language is null) and d.tenant = ? and d.key = ? and d.type = ?`

	base, _ := lang.Base()

	slog.Debugf("loading for language %s", base.String())

	_, err := s.db.QueryOne(&dict, sql, base.String(), tenant, key, dictionaryType)
	return &dict, err

}

func (s *dictionaryRepository) Save(dict *model.Dictionary) error {
	_, err := s.db.Model(dict).Insert()
	return err
}

func (s *dictionaryRepository) Update(dict *model.Dictionary) error {
	res, err := s.db.Model(dict).WherePK().Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("not existing key=%s tenant=%s and type=%s", dict.Key, dict.Tenant, dict.Type)
	}

	return nil
}

func (s *dictionaryRepository) LoadAll(tenant string) ([]model.Dictionary, error) {

	var dicts []model.Dictionary

	_, err := s.db.Query(&dicts,
		`select * from all_dictionaries where tenant = ?`, tenant)

	return dicts, err
}

func (s *dictionaryRepository) LoadChildren(parentKey, dictionaryType, tenant string) ([]model.Dictionary, error) {
	var children []model.Dictionary

	sql := `select * from dictionary where parent_key = ? and type = ? and tenant = ?`

	_, err := s.db.Query(&children, sql, parentKey, dictionaryType, tenant)

	return children, err
}

func (s *dictionaryRepository) LoadChildrenDictionaryTranslated(parentKey, dictionaryType, tenant string, lang language.Tag) ([]model.ChildDictionary, error) {

	var children []model.ChildDictionary

	res, err := s.LoadChildrenTranslated(parentKey, dictionaryType, tenant, lang)

	if err != nil {
		return nil, err
	}

	for _, item := range res {
		children = append(children, model.ChildDictionary{
			Key:       item.Key,
			Name:      item.Name,
			ParentKey: *item.ParentKey,
			Content:   item.Content,
		})
	}
	return children, nil
}

func (s *dictionaryRepository) LoadChildrenTranslated(parentKey, dictionaryType, tenant string, lang language.Tag) ([]model.Dictionary, error) {
	var children []model.Dictionary

	sql := `select d.key, d.type, COALESCE(t.name, d.name) as name , group_id, d.tenant, content, parent_key
			from dictionary d
				left outer join translation t on d.key = t.key and d.tenant = t.tenant and d.type = t.type
			where (t.language = ? or t.language is null) and parent_key = ? and d.type = ? and d.tenant = ?`

	base, _ := lang.Base()
	slog.Debugf("loading for language %s", base.String())

	_, err := s.db.Query(&children, sql, base.String(), parentKey, dictionaryType, tenant)

	return children, err
}

type ChildrenKeys map[string]bool

func (s *dictionaryRepository) LoadChildrenKeys(parentKey, dictionaryType, tenant string) (ChildrenKeys, error) {

	var keys []string
	_, err := s.db.Query(&keys, `select key from dictionary where parent_key = ? and type = ? and tenant = ?`,
		parentKey, dictionaryType, tenant)

	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	for _, k := range keys {
		res[k] = true
	}

	return res, nil
}

func (s *dictionaryRepository) LoadByType(dictionaryType, tenant string) ([]model.Dictionary, error) {

	var res []model.Dictionary
	err := s.db.Model(&res).Where("type = ? and tenant = ? and parent_key is null", dictionaryType, tenant).
		Select()
	return res, err
}

func (s *dictionaryRepository) DeleteMultiple(keys ChildrenKeys, tenant, dictionaryType string) error {

	var toDelete []model.Dictionary

	for k := range keys {
		toDelete = append(toDelete, model.Dictionary{
			Key:    k,
			Type:   dictionaryType,
			Tenant: tenant,
		})
	}

	_, err := s.db.Model(&toDelete).WherePK().Delete()

	return err
}

func (s *dictionaryRepository) Delete(key, dictionaryType, tenant string) error {
	_, err := s.db.Model(&model.Dictionary{Key: key, Type: dictionaryType, Tenant: tenant}).
		WherePK().
		Delete()
	return err
}

func (s *dictionaryRepository) DeleteAll(tenant string) error {
	_, err := s.db.Model(&model.Dictionary{Tenant: tenant}).Where("tenant = ?tenant").Delete()
	return err
}

func (s *dictionaryRepository) DeleteByType(dictionaryType, tenant string) error {
	_, err := s.db.Model(&model.Dictionary{Tenant: tenant, Type: dictionaryType}).
		Where("tenant = ?tenant and type = ?type").Delete()
	return err
}
