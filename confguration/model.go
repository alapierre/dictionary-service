package confguration

import (
	"database/sql"
	"time"
)

type Configuration struct {
	tableName struct{} `pg:"configuration,alias:c"`

	Key      string         `pg:"key,pk" json:"key"`
	Tenant   string         `pg:"tenant,pk" json:"tenant"`
	Type     string         `json:"type"`
	Name     string         `json:"name"`
	Value    sql.NullString `json:"value"`
	DateFrom time.Time      `pg:"date_from,pk" json:"date_from"`
	DateTo   time.Time      `json:"date_to"`
}
