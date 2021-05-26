package calendar

import (
	"time"
)

type DictionaryCalendar struct {
	tableName struct{} `pg:"calendar,alias:c,discard_unknown_columns"`

	Day    time.Time         `pg:"day,pk" `
	Tenant string            `pg:"tenant,pk"`
	Name   *string           `pg:"name"`
	Type   string            `pg:"type,use_zero"`
	Kind   *string           `pg:"kind"`
	Labels map[string]string `pg:"labels,hstore"`
}

type DictionaryCalendarType struct {
	tableName struct{} `pg:"calendar_type,alias:t,discard_unknown_columns"`

	Type   string `pg:"type,pk"`
	Name   string `pg:"name"`
	Tenant string `pg:"tenant,pk"`
}
