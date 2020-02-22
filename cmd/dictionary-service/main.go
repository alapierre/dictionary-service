package main

import (
	"dictionaries-service/service"
	"dictionaries-service/util"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v9"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerPort           int    `split_words:"true" default:"9098"`
	DatasourceName       string `split_words:"true" default:"app"`
	DatasourceSchema     string `split_words:"true" default:"dictionary"`
	DatasourceHost       string `split_words:"true" default:"localhost:5432"`
	DatasourceUser       string `required:"true" split_words:"true"`
	DatasourcePassword   string `required:"true" split_words:"true"`
	DatasourcePoolSize   int    `split_words:"true" default:"2"`
	DatasourceMaxRetries int    `split_words:"true" default:"3"`
	EurekaServiceUrl     string `split_words:"true" default:"http://localhost:8761/eureka"`
}

var c Config

func main() {

	slog.Info("Starting up Dictionary Service")
	err := envconfig.Process("dict", &c)
	util.FailOnError(err, "Can't parse environment variables")

	db := connectDb()
	defer util.Close(db)

	repository := service.NewDictionaryRepository(db)
	dicts, err := repository.LoadAll("")

	util.FailOnError(err, "Can't query")

	for _, d := range dicts {
		fmt.Printf("%#v\n", d)
	}

}

func connectDb() *pg.DB {

	params := make(map[string]interface{})
	params["search_path"] = "dictionary"

	db := pg.Connect(&pg.Options{
		User:       c.DatasourceUser,
		Password:   c.DatasourcePassword,
		Addr:       c.DatasourceHost,
		Database:   c.DatasourceName,
		PoolSize:   c.DatasourcePoolSize,
		MaxRetries: c.DatasourceMaxRetries,

		//OnConnect: func(conn *pg.Conn) error {
		//	_, err := conn.Exec("set search_path=?", "scheduler")
		//	if err != nil {
		//		slog.Error(err)
		//	}
		//	return nil
		//},
	})

	//err := model.CreateSchema(db)
	//util.FailOnError(err, "Cant create schema")

	return db
}
