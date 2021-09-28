package configuration

import (
	"dictionaries-service/util"
	test2 "dictionaries-service/util/test"
	slog "github.com/go-eden/slf4go"
	"os"
	"testing"
)

var confRepository Repository

func TestMain(m *testing.M) {

	db := test2.ConnectDb()

	confRepository = NewRepository(db)

	ex := m.Run()

	slog.Info("Shutting down test init func")
	util.Close(db)

	os.Exit(ex)

}
