package calendar

import (
	"dictionaries-service/types"
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

	Type   string `pg:"type,pk" json:"type"`
	Name   string `pg:"name" json:"name"`
	Tenant string `pg:"tenant,pk,use_zero" json:"-"`
}

type SaveDto struct {
	Day          types.JsonDate    `json:"day"`
	CalendarType string            `json:"-"`
	Name         string            `json:"name"`
	Kind         *string           `json:"kind,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}
