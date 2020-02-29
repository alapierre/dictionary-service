package service

import (
	"context"
	"dictionaries-service/util"
	"fmt"
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
	db.AddQueryHook(dbLogger{})
	return db
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Print("go-pg: ")
	fmt.Println(q.FormattedQuery())
	return nil
}
