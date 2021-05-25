package calendar

import (
	"github.com/go-pg/pg/v10"
	"time"
)

func NewRepository(db *pg.DB) Repository {
	return &calendarRepository{db: db}
}

type Repository interface {
	LoadByTypeAndRange(tenant string, calendarType string, from, to time.Time) ([]DictionaryCalendar, error)
}

type calendarRepository struct {
	db *pg.DB
}

func (c *calendarRepository) LoadByTypeAndRange(tenant string, calendarType string, from, to time.Time) ([]DictionaryCalendar, error) {

	var result []DictionaryCalendar

	err := c.db.Model(&result).
		Where("type = ? and tenant = ? and day >= ? and day <= ?", calendarType, tenant, from, to).
		Select()

	return result, err
}
