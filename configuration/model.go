package configuration

import (
	"database/sql"
	"time"
)

type Short struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Configuration struct {
	tableName struct{} `pg:"configuration,alias:c"`

	Key      string         `pg:"key,pk" json:"key"`
	Tenant   string         `pg:"tenant,pk,use_zero" json:"tenant"`
	Type     string         `json:"type"`
	Name     string         `json:"name"`
	Value    sql.NullString `json:"value"`
	DateFrom time.Time      `pg:"date_from,pk" json:"date_from"`
	DateTo   time.Time      `json:"date_to"`
}
