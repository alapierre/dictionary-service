package calendar

import (
	"dictionaries-service/util"
	"fmt"
	"github.com/go-pg/pg/v10"
)

func NewTypeRepository(db *pg.DB) TypeRepository {
	return &calendarTypeRepository{db: db}
}

type TypeRepository interface {
	LoadAll(tenant string) ([]DictionaryCalendarType, error)
	Save(cal *DictionaryCalendarType) error
	Update(cal *DictionaryCalendarType) error
	Delete(tenant string, calendarType string) error
	Load(tenant, calendarType string) (*DictionaryCalendarType, error)
}

type calendarTypeRepository struct {
	db *pg.DB
}

func (c calendarTypeRepository) Load(tenant, calendarType string) (*DictionaryCalendarType, error) {

	result := &DictionaryCalendarType{Tenant: tenant, Type: calendarType}

	err := c.db.Model(result).
		WherePK().
		Select()

	return result, err
}

func (c calendarTypeRepository) LoadAll(tenant string) ([]DictionaryCalendarType, error) {

	var result []DictionaryCalendarType

	err := c.db.Model(&result).
		Where("tenant = ?", tenant).
		Select()

	return result, err
}

func (c calendarTypeRepository) Save(cal *DictionaryCalendarType) error {
	_, err := c.db.Model(cal).Insert()
	return err
}

func (c calendarTypeRepository) Update(cal *DictionaryCalendarType) error {

	res, err := c.db.Model(cal).WherePK().Update()
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return util.NewZeroRowsAffectedError(fmt.Errorf("there is no row for geven key %s and tenant %s", cal.Type, cal.Tenant))
	}

	return nil
}

func (c calendarTypeRepository) Delete(tenant string, calendarType string) error {

	res, err := c.db.Model(&DictionaryCalendarType{Tenant: tenant, Type: calendarType}).
		WherePK().
		Delete()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return util.NewZeroRowsAffectedError(fmt.Errorf("can't delete no existent row for geven key %s and tenant %s", calendarType, tenant))
	}

	return nil
}
