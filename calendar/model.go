package calendar

import (
	"time"
)

type DictionaryCalendar struct {
	tableName struct{} `pg:"calendar,alias:c,discard_unknown_columns"`

	Day    time.Time         `pg:"day,pk" `
	Tenant string            `pg:"tenant,pk,use_zero"`
	Name   *string           `pg:"name"`
	Type   string            `pg:"type,pk,use_zero"`
	Kind   *string           `pg:"kind"`
	Labels map[string]string `pg:"labels,hstore"`
}

type DictionaryCalendarType struct {
	tableName struct{} `pg:"calendar_type,alias:t,discard_unknown_columns"`

	Type   string `pg:"type,pk"`
	Name   string `pg:"name"`
	Tenant string `pg:"tenant,pk,use_zero"`
}

type SaveDto struct {
	Day          time.Time         `json:"-"`
	CalendarType string            `json:"-"`
	Name         string            `json:"name"`
	Kind         *string           `json:"kind,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}
