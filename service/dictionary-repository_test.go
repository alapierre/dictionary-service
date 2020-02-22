package service

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"testing"
)

var repo DictionaryRepository

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

}

func Test_dictionaryRepository_Load(t *testing.T) {

	dict, err := repo.Load("uw", "AbsenceType", "")

	if err != nil {
		t.Fatal(err, "Can't query")
	}

	fmt.Printf("%#v\n", dict)
}

func Test_dictionaryRepository_LoadAll(t *testing.T) {

	dicts, err := repo.LoadAll("")

	if err != nil {
		t.Fatal(err, "Can't query")
	}

	for _, d := range dicts {
		fmt.Printf("%#v\n", d)
	}

}

func Test_dictionaryRepository_LoadByType(t *testing.T) {

}

func Test_dictionaryRepository_Save(t *testing.T) {

}
