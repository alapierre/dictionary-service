package service

import (
	"dictionaries-service/util"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v9"
	"os"
	"testing"
)

var (
	db             *pg.DB
	service        *DictionaryService
	dictRepository DictionaryRepository
)

func TestMain(m *testing.M) {

	slog.Info("Running test init func")

	db = connectTestDb()

	dictRepository = NewDictionaryRepository(db)
	service = NewDictionaryService(dictRepository)

	ex := m.Run()

	slog.Info("Shutting down test init func")
	util.Close(db)

	os.Exit(ex)
}

func connectTestDb() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "app",
		Password: "qwedsazxc",
		Addr:     "localhost:5432",
		Database: "app",
		//OnConnect: func(conn *pg.Conn) error {
		//	_, err := conn.Exec("set search_path=?", "scheduler")
		//	if err != nil {
		//		slog.Error(err)
		//	}
		//	return nil
		//},
	})
	return db
}
