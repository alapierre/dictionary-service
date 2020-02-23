package main

import (
	"context"
	"dictionaries-service/service"
	"dictionaries-service/transport"
	"dictionaries-service/util"
	"flag"
	"fmt"
	slog "github.com/go-eden/slf4go"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	slog.Infof("database name: %s host: %s user: %s", c.DatasourceName, c.DatasourceHost, c.DatasourceUser)

	db := connectDb()
	defer util.Close(db)

	//migrate(db)

	dictRepository := service.NewDictionaryRepository(db)
	chRepository := service.NewChildRepository(db)
	service := service.NewDictionaryService(dictRepository, chRepository)

	r := mux.NewRouter()
	r.Use(addContext)

	r.Methods("GET").Path("/api/dictionary/{type}/{key}").Handler(httptransport.NewServer(
		transport.MakeLoadDictEndpoint(service),
		transport.LoadDictRequest,
		transport.EncodeResponse,
	))

	http.Handle("/", r)

	slog.Info("Started on port ", c.ServerPort)

	startHttpAndWaitForSigINT(c.ServerPort)

	slog.Info("Bye.")
}

func startHttpAndWaitForSigINT(port int) {

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				slog.Panicf("can't start HTTP server at port %d, %v", port, err)
			}
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	slog.Info("Shutdown in progres")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Fatalf("Can't shutdown HTTP server %v", err)
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

func migrate(db *pg.DB) {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("0 args ")
		_, _, err := migrations.Run(db, "init")
		if err != nil {
			slog.Infof("initializing migration %v", err)
		}
	}

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	util.FailOnError(err, "Can't run migrations")

	if newVersion != oldVersion {
		slog.Infof("db schema migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		slog.Infof("db schema version is %d, migration is not necessary\n", oldVersion)
	}
}

func addContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "tenant", "")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
