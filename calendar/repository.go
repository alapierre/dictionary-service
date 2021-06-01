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
	Save(cal *DictionaryCalendar) error
	Update(cal *DictionaryCalendar) error
	Delete(tenant string, calendarType string, day time.Time) error
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

func (c *calendarRepository) Save(cal *DictionaryCalendar) error {
	_, err := c.db.Model(cal).Insert()
	return err
}

func (c *calendarRepository) Update(cal *DictionaryCalendar) error {
	_, err := c.db.Model(cal).WherePK().Update()
	return err
}

func (c *calendarRepository) Delete(tenant string, calendarType string, day time.Time) error {
	_, err := c.db.Model(&DictionaryCalendar{Day: day, Tenant: tenant, Type: calendarType}).
		WherePK().
		Delete()
	return err
}
