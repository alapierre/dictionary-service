package calendar

import (
	"dictionaries-service/util"
	test2 "dictionaries-service/util/test"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	db              *pg.DB
	repository      Repository
	calendarService Service
)

func Test_calendarRepository_Load(t *testing.T) {

	var from = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	var to = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

	cal, err := repository.LoadByTypeAndRange("", "holiday", from, to)

	assert.NoError(t, err)

	fmt.Printf("res: %#v\n", cal)
}

func TestMain(m *testing.M) {

	slog.Info("Running test init func")

	db = test2.ConnectDb()
	repository = NewRepository(db)
	calendarService = NewService(repository, NewTypeRepository(db))

	ex := m.Run()

	slog.Info("Shutting down test init func")
	util.Close(db)

	os.Exit(ex)
}
