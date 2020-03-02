package service

import (
	"dictionaries-service/model"
	"dictionaries-service/util"
	"github.com/go-pg/pg/v9"
)

func NewTranslateRepository(db *pg.DB) TranslateRepository {
	return &translateRepository{db: db}
}

type TranslateRepository interface {
	DoInTransaction(callback util.TransactionCallback) error
	Save(translate *model.Translation) error
	Update(translate *model.Translation) error
	LoadByTypeAndLang(dictionaryType, lang, tenant string) ([]model.Translation, error)
	DeleteByKeyAndType(key, dictionaryType, tenant string) error
	Delete(translation *model.Translation) error
	AvailableTranslation(key, dictionaryType, tenant string) ([]string, error)
}

type Translator interface {
	AvailableTranslation(key, dictionaryType, tenant string) ([]string, error)
}

type translateRepository struct {
	db *pg.DB
}

func (s *translateRepository) DoInTransaction(callback util.TransactionCallback) error {
	return util.DoInTransaction(s.db, callback)
}

func (s *translateRepository) Save(translate *model.Translation) error {
	err := s.db.Insert(translate)
	return err
}

func (s *translateRepository) Update(translate *model.Translation) error {
	err := s.db.Update(translate)
	return err
}

func (s *translateRepository) LoadByTypeAndLang(dictionaryType, lang, tenant string) ([]model.Translation, error) {

	var translations []model.Translation

	_, err := s.db.Query(&translations,
		`select * from translation where tenant = ? and type = ? and  language = ?`, tenant, dictionaryType, lang)

	return translations, err
}

func (s *translateRepository) DeleteByKeyAndType(key, dictionaryType, tenant string) error {
	_, err := s.db.Model(&model.Translation{Tenant: tenant, Type: dictionaryType, Key: key}).
		Where("tenant = ?tenant and type = ?type and key = ?key").Delete()
	return err
}

func (s *translateRepository) Delete(translation *model.Translation) error {
	return s.db.Delete(translation)
}

func (s *translateRepository) AvailableTranslation(key, dictionaryType, tenant string) ([]string, error) {

	type langName struct {
		Language string `pg:"language"`
	}

	var translations []langName

	sql := `select language from translation where key = ? and type = ? and tenant = ?`

	_, err := s.db.Query(&translations, sql, key, dictionaryType, tenant)

	if err != nil {
		return nil, err
	}

	var res []string

	for _, t := range translations {
		res = append(res, t.Language)
	}

	return res, nil
}
