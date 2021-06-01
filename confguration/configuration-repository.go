package confguration

import (
	"github.com/go-pg/pg/v10"
	"time"
)

var EndOfTheWorld = time.Date(2999, 12, 31, 0, 0, 0, 0, time.UTC)
var BeginOfTheWorld = time.Date(1900, 12, 31, 0, 0, 0, 0, time.UTC)

func NewRepository(db *pg.DB) Repository {
	return &repository{db: db}
}

type Repository interface {
	LoadForDay(key, tenant string, day time.Time) (*Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]Configuration, error)
	Save(configuration *Configuration) error
}

type repository struct {
	db *pg.DB
}

func (c *repository) LoadForDay(key, tenant string, day time.Time) (*Configuration, error) {

	config := &Configuration{}
	err := c.db.Model(config).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, day, day).
		Select()

	return config, err
}

func (c *repository) Load(key, tenant string, from, to time.Time) ([]Configuration, error) {

	var result []Configuration
	err := c.db.Model(&result).
		Where("key = ? and tenant = ? and date_from <= ? and date_to >= ?", key, tenant, to, from).
		Select()

	return result, err
}

func (c *repository) Save(configuration *Configuration) error {
	_, err := c.db.Model(configuration).Insert()
	return err
}
