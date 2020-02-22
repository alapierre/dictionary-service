package service

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"testing"
)

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

	chRepository = NewChildRepository(db)
}

func TestChildRepository_LoadChildren(t *testing.T) {
	ch, err := chRepository.LoadChildren("uw", "AbsenceType", "")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", ch)

}
