package service

import (
	"dictionaries-service/model"
	"github.com/go-pg/pg/v10"
	"time"
)

var EndOfTheWorld = time.Date(2999, 12, 31, 0, 0, 0, 0, time.UTC)
var BeginOfTheWorld = time.Date(1900, 12, 31, 0, 0, 0, 0, time.UTC)

func NewConfigurationRepository(db *pg.DB) ConfigurationRepository {
	return &configurationRepository{db: db}
}

type ConfigurationRepository interface {
	LoadForDay(key, tenant string, day time.Time) (*model.Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]model.Configuration, error)
	Save(configuration *model.Configuration) error
}

type configurationRepository struct {
	db *pg.DB
}

func (c *configurationRepository) LoadForDay(key, tenant string, day time.Time) (*model.Configuration, error) {

	config := &model.Configuration{}
	err := c.db.Model(config).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, day, day).
		Select()

	return config, err
}

func (c *configurationRepository) Load(key, tenant string, from, to time.Time) ([]model.Configuration, error) {

	var result []model.Configuration
	err := c.db.Model(&result).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, to, from).
		Select()

	return result, err
}

func (c *configurationRepository) Save(configuration *model.Configuration) error {
	_, err := c.db.Model(configuration).Insert()
	return err
}
