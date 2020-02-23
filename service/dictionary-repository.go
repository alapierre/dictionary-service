package service

import (
	"dictionaries-service/model"
	"github.com/go-pg/pg/v9"
)

func NewDictionaryRepository(db *pg.DB) DictionaryRepository {
	return &dictionaryRepository{db: db}
}

type DictionaryRepository interface {
	Load(key, dictionaryType, tenant string) (*model.Dictionary, error)
	Save(dict *model.Dictionary) error
	LoadAll(tenant string) ([]model.Dictionary, error)
	LoadByType(dictionaryType, tenant string) ([]model.Dictionary, error)
}

type dictionaryRepository struct {
	db *pg.DB
}

func (s *dictionaryRepository) Load(key, dictionaryType, tenant string) (*model.Dictionary, error) {

	var dict model.Dictionary
	_, err := s.db.QueryOne(&dict,
		`select * from all_dictionaries where tenant = ? and key = ? and type = ?`, tenant, key, dictionaryType)

	return &dict, err

}

func (s *dictionaryRepository) Save(dict *model.Dictionary) error {
	return nil
}

func (s *dictionaryRepository) LoadAll(tenant string) ([]model.Dictionary, error) {

	var dicts []model.Dictionary

	_, err := s.db.Query(&dicts,
		`select * from all_dictionaries where tenant = ?`, tenant)

	return dicts, err
}

func (s *dictionaryRepository) LoadByType(dictionaryType, tenant string) ([]model.Dictionary, error) {
	return nil, nil
}
