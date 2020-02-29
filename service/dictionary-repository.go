package service

import (
	"dictionaries-service/model"
	"dictionaries-service/util"
	"github.com/go-pg/pg/v9"
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
	DeleteMultiple(keys ChildrenKeys) error
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

func (s *dictionaryRepository) Save(dict *model.Dictionary) error {
	err := s.db.Insert(dict)
	return err
}

func (s *dictionaryRepository) Update(dict *model.Dictionary) error {
	err := s.db.Update(dict)
	return err
}

func (s *dictionaryRepository) LoadAll(tenant string) ([]model.Dictionary, error) {

	var dicts []model.Dictionary

	_, err := s.db.Query(&dicts,
		`select * from all_dictionaries where tenant = ?`, tenant)

	return dicts, err
}

func (s *dictionaryRepository) LoadChildren(parentKey, dictionaryType, tenant string) ([]model.Dictionary, error) {
	var children []model.Dictionary
	_, err := s.db.Query(&children, `select * from dictionary where parent_key = ? and type = ? and tenant = ?`,
		parentKey, dictionaryType, tenant)

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
	return nil, nil
}

func (s *dictionaryRepository) DeleteMultiple(keys ChildrenKeys) error {

	var toDelete []model.Dictionary

	for k := range keys {
		toDelete = append(toDelete, model.Dictionary{
			Key: k,
		})
	}

	_, err := s.db.Model(&toDelete).WherePK().Delete()

	return err
}
