package service

import (
	"dictionaries-service/model"
	"dictionaries-service/util"
	"github.com/go-pg/pg/v9"
)

func NewDictionaryMetadataRepository(db *pg.DB) DictionaryMetadataRepository {
	return &dictionaryMetadataRepository{db: db}
}

type DictionaryMetadataRepository interface {
	DoInTransaction(callback util.TransactionCallback) error
	Save(translate *model.DictionaryMetadata) error
	Update(translate *model.DictionaryMetadata) error
	Delete(translation *model.DictionaryMetadata) error
	Load(dictionaryType, tenant string) (*model.DictionaryMetadata, error)
	AvailableDictionaryTypes(tenant string) ([]string, error)
}

type dictionaryMetadataRepository struct {
	db *pg.DB
}

func (s *dictionaryMetadataRepository) Load(dictionaryType, tenant string) (*model.DictionaryMetadata, error) {

	metadata := &model.DictionaryMetadata{
		Type:   dictionaryType,
		Tenant: tenant,
	}
	err := s.db.Select(metadata)

	return metadata, err
}

func (s *dictionaryMetadataRepository) DoInTransaction(callback util.TransactionCallback) error {
	return util.DoInTransaction(s.db, callback)
}

func (s *dictionaryMetadataRepository) Save(metadata *model.DictionaryMetadata) error {
	err := s.db.Insert(metadata)
	return err
}

func (s *dictionaryMetadataRepository) Update(metadata *model.DictionaryMetadata) error {
	err := s.db.Update(metadata)
	return err
}

func (s *dictionaryMetadataRepository) Delete(metadata *model.DictionaryMetadata) error {
	return s.db.Delete(metadata)
}

func (s *dictionaryMetadataRepository) AvailableDictionaryTypes(tenant string) ([]string, error) {

	var types []string
	_, err := s.db.Query(&types, `select type from dictionary_metadata dictionary where tenant = ?`, tenant)

	return types, err
}
