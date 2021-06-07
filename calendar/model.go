package calendar

import (
	"time"
)

// DictionaryCalendar calendar item
// swagger:model DictionaryCalendar
type DictionaryCalendar struct {
	tableName struct{} `pg:"calendar,alias:c,discard_unknown_columns"`

	Day    time.Time         `pg:"day,pk" `
	Tenant string            `pg:"tenant,pk,use_zero"`
	Name   *string           `pg:"name"`
	Type   string            `pg:"type,pk,use_zero"`
	Kind   *string           `pg:"kind"`
	Labels map[string]string `pg:"labels,hstore"`
}

// DictionaryCalendarType calendar type
// swagger:model DictionaryCalendarType
type DictionaryCalendarType struct {
	tableName struct{} `pg:"calendar_type,alias:t,discard_unknown_columns"`

	Type   string `pg:"type,pk"`
	Name   string `pg:"name"`
	Tenant string `pg:"tenant,pk,use_zero"`
}
