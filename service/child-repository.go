package service

import (
	"dictionaries-service/model"
	"github.com/go-pg/pg/v9"
)

func NewChildRepository(db *pg.DB) ChildRepository {
	return &childRepository{db: db}
}

type ChildRepository interface {
	LoadChildren(parentKey, dictionaryType, tenant string) ([]model.Child, error)
	Load(Key, dictionaryType, tenant string) (*model.Child, error)
	Save(dict *model.Child) error
	Update(dict *model.Child) error
	Delete(key, dictionaryType, tenant string) error
}

type childRepository struct {
	db *pg.DB
}

func (r *childRepository) LoadChildren(parentKey, dictionaryType, tenant string) ([]model.Child, error) {
	var children []model.Child
	err := r.db.Model(&children).Where("parent_key = ? and type = ? and tenant = ?",
		parentKey, dictionaryType, tenant).
		Select()
	return children, err
}

func (r *childRepository) Load(key, dictionaryType, tenant string) (*model.Child, error) {
	var child = &model.Child{
		Key:    key,
		Type:   dictionaryType,
		Tenant: tenant,
	}
	err := r.db.Select(child)
	return child, err
}

func (r *childRepository) Save(dict *model.Child) error {
	err := r.db.Insert(dict)
	return err
}

func (r *childRepository) Update(dict *model.Child) error {
	err := r.db.Update(dict)
	return err
}

func (r *childRepository) Delete(key, dictionaryType, tenant string) error {
	return r.db.Delete(&model.Child{
		Key:    key,
		Tenant: tenant,
		Type:   dictionaryType,
	})
}
