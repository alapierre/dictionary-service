package service

import (
	"dictionaries-service/model"
	"dictionaries-service/util"
	"github.com/go-pg/pg/v10"
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
	err := s.db.Model(metadata).
		WherePK().
		Select()

	return metadata, err
}

func (s *dictionaryMetadataRepository) DoInTransaction(callback util.TransactionCallback) error {
	return util.DoInTransaction(s.db, callback)
}

func (s *dictionaryMetadataRepository) Save(metadata *model.DictionaryMetadata) error {
	_, err := s.db.Model(metadata).Insert()
	return err
}

func (s *dictionaryMetadataRepository) Update(metadata *model.DictionaryMetadata) error {
	_, err := s.db.Model(metadata).WherePK().Update()
	return err
}

func (s *dictionaryMetadataRepository) Delete(metadata *model.DictionaryMetadata) error {
	_, err := s.db.Model(metadata).WherePK().Delete()
	return err
}

func (s *dictionaryMetadataRepository) AvailableDictionaryTypes(tenant string) ([]string, error) {

	var types []string
	_, err := s.db.Query(&types, `select type from dictionary_metadata where tenant = ?`, tenant)

	return types, err
}
