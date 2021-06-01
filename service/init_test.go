package service

import (
	"dictionaries-service/test"
	"dictionaries-service/util"
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

	db = test.ConnectDb()

	dictRepository = NewDictionaryRepository(db)
	translationRepository := NewTranslateRepository(db)
	metadataRepository = NewDictionaryMetadataRepository(db)
	service = NewDictionaryService(dictRepository, translationRepository, metadataRepository)

	ex := m.Run()

	slog.Info("Shutting down test init func")
	util.Close(db)

	os.Exit(ex)
}
