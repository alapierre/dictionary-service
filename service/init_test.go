package service

import (
	"context"
	"dictionaries-service/util"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v10"
	"os"
	"testing"
)

var (
	db                 *pg.DB
	service            *DictionaryService
	dictRepository     DictionaryRepository
	metadataRepository DictionaryMetadataRepository
)

func TestMain(m *testing.M) {

	slog.Info("Running test init func")

	db = connectTestDb()

	dictRepository = NewDictionaryRepository(db)
	translationRepository := NewTranslateRepository(db)
	metadataRepository = NewDictionaryMetadataRepository(db)
	service = NewDictionaryService(dictRepository, translationRepository, metadataRepository)

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

func (d dbLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	fmt.Print("go-pg: ")
	fmt.Println(q.FormattedQuery())
	return nil
}
