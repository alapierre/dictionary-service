package util

import (
	"github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v9"
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
