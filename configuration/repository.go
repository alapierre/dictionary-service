package configuration

import (
	"dictionaries-service/util"
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"
)

var EndOfTheWorld = time.Date(2999, 12, 31, 0, 0, 0, 0, time.UTC)
var BeginOfTheWorld = time.Date(1900, 12, 31, 0, 0, 0, 0, time.UTC)

func NewRepository(db *pg.DB) Repository {
	return &configurationRepository{db: db}
}

type Repository interface {
	LoadForDay(key, tenant string, day time.Time) (*Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]Configuration, error)
	Save(configuration *Configuration) error
	LoadAllShort(tenant string) ([]Short, error)
	LoadValues(tenant, key string) ([]Configuration, error)
	Update(configuration *Configuration) error
	DeleteValue(tenant, key string, dateFrom time.Time) error
	LoadFirst(key, tenant string) (*Configuration, error)
	LoadById(key, tenant string, from time.Time) (Configuration, error)
}

type configurationRepository struct {
	db *pg.DB
}

func (c *configurationRepository) LoadForDay(key, tenant string, day time.Time) (*Configuration, error) {

	config := &Configuration{}
	err := c.db.Model(config).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, day, day).
		Select()

	return config, err
}

func (c *configurationRepository) LoadFirst(key, tenant string) (*Configuration, error) {

	config := &Configuration{}
	_, err := c.db.Query(config, `select * from dictionary.configuration where key = ? and tenant = ? order by date_from limit 1`, key, tenant)

	return config, err

}

func (c *configurationRepository) LoadById(key, tenant string, from time.Time) (Configuration, error) {

	var result = Configuration{
		Key:      key,
		Tenant:   tenant,
		DateFrom: from,
	}

	err := c.db.Model(&result).
		WherePK().
		Select()

	return result, err
}

func (c *configurationRepository) Load(key, tenant string, from, to time.Time) ([]Configuration, error) {

	var result []Configuration
	err := c.db.Model(&result).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, to, from).
		Select()

	return result, err
}

func (c *configurationRepository) Save(configuration *Configuration) error {
	_, err := c.db.Model(configuration).Insert()
	return err
}

func (c *configurationRepository) Update(configuration *Configuration) error {

	res, err := c.db.Model(configuration).WherePK().Update()

	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return util.NewZeroRowsAffectedError(fmt.Errorf("there is no row for geven dateFrom: %s, type: %s and tenant: %s",
			configuration.DateFrom, configuration.Key, configuration.Tenant))
	}
	return nil
}

func (c *configurationRepository) LoadAllShort(tenant string) ([]Short, error) {
	var result []Short
	_, err := c.db.Query(&result, `select distinct key, name, type from configuration where tenant = ?`, tenant)
	return result, err
}

func (c *configurationRepository) LoadValues(tenant, key string) ([]Configuration, error) {
	var result []Configuration
	err := c.db.Model(&result).
		Where("key = ? and tenant = ?", key, tenant).
		Select()

	return result, err
}

func (c *configurationRepository) DeleteValue(tenant, key string, dateFrom time.Time) error {

	res, err := c.db.Model(&Configuration{
		Key:      key,
		Tenant:   tenant,
		DateFrom: dateFrom,
	}).
		WherePK().
		Delete()

	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return util.NewZeroRowsAffectedError(fmt.Errorf("there is no row for geven day: %s, key: %s and tenant: %s", dateFrom, key, tenant))
	}
	return nil
}
