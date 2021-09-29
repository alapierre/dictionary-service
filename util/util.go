package util

import (
	"database/sql"
	"github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v10"
	"time"
)

func FailOnError(err error, msg string) {
	if err != nil {
		slog.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}

type Closable interface {
	Close() error
}

func Close(target Closable) {
	err := target.Close()
	if err != nil {
		slog.Errorf("Cant close %v", target)
	}
}

type TransactionCallback func() error

func DoInTransaction(db *pg.DB, callback TransactionCallback) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = callback()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func SqlNullStringToStringPointer(str sql.NullString) *string {
	if str.Valid {
		tmp := str.String
		return &tmp
	} else {
		return nil
	}
}

func PointerToSqlNullString(str *string) sql.NullString {

	if str == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{
		String: *str,
		Valid:  true,
	}
}

func StringToTime(str string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	return time.Parse(layout, str)
}
