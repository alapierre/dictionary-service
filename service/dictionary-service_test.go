package service

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"testing"
)

var chRepository ChildRepository
var service *DictionaryService

func init() {

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

	repo = NewDictionaryRepository(db)
	chRepository = NewChildRepository(db)
	service = NewDictionaryService(repo, chRepository)

}

func TestDictionaryService_Load(t *testing.T) {

}

func TestDictionaryService_LoadShallow(t *testing.T) {

	dict, err := service.Load("uw", "AbsenceType", "")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", dict)

}

func TestNewDictionaryService(t *testing.T) {

}

func Test_mergeMaps(t *testing.T) {

}

func Test_prepareChildrenMap(t *testing.T) {

}

func Test_prepareMap(t *testing.T) {

}
