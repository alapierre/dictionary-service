//go:generate swagger generate spec -o ./../../swagger.yml
package main

import (
	"context"
	"dictionaries-service/calendar"
	"dictionaries-service/configuration"
	"dictionaries-service/service"
	"dictionaries-service/tenant"
	rest "dictionaries-service/transport/http"
	"dictionaries-service/util"
	"flag"
	"fmt"
	"github.com/alapierre/gokit-utils/eureka"
	slog "github.com/go-eden/slf4go"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	ServerPort               int    `split_words:"true" default:"9098"`
	DatasourceName           string `split_words:"true" default:"app"`
	DatasourceSchema         string `split_words:"true" default:"dictionary"`
	DatasourceHost           string `split_words:"true" default:"localhost:5432"`
	DatasourceUser           string `required:"true" split_words:"true"`
	DatasourcePassword       string `required:"true" split_words:"true"`
	DatasourcePoolSize       int    `split_words:"true" default:"2"`
	DatasourceMaxRetries     int    `split_words:"true" default:"3"`
	EurekaServiceUrl         string `split_words:"true" default:"http://localhost:8761/eureka"`
	DefaultLanguage          string `split_words:"true" default:"en"`
	InitDBConnectionRets     int    `split_words:"true" default:"100"`
	ShowSql                  bool   `split_words:"true" default:"true"`
	TenantHeaderName         string `split_words:"true" default:"X-Tenant-ID"`
	EnableEureka             bool   `split_words:"true" default:"true"`
	AccessControlAllowOrigin string `split_words:"true" default:"*"`
}

var c Config

func main() {

	slog.Info("Starting up Dictionary Service")
	err := envconfig.Process("dict", &c)
	util.FailOnError(err, "Can't parse environment variables")

	slog.Infof("database name: %s host: %s user: %s", c.DatasourceName, c.DatasourceHost, c.DatasourceUser)

	rest.DefaultLanguage = language.MustParse(c.DefaultLanguage)
	slog.Info("Default language is ", rest.DefaultLanguage)

	tenant.HeaderName(c.TenantHeaderName, "X-Tenant-Merge-Default")
	slog.Info("Tenant header is ", c.TenantHeaderName)

	db := connectDb()
	defer util.Close(db)

	migrate(db)

	dictionaryRepository := service.NewDictionaryRepository(db)
	translationRepository := service.NewTranslateRepository(db)
	metadataRepository := service.NewDictionaryMetadataRepository(db)
	dictionaryService := service.NewDictionaryService(dictionaryRepository, translationRepository, metadataRepository)

	configurationRepository := configuration.NewRepository(db)
	configurationService := configuration.NewConfigurationService(configurationRepository)

	calendarService := calendar.NewService(calendar.NewRepository(db), calendar.NewTypeRepository(db))

	r := mux.NewRouter()
	r.Use(addContext, accessControlMiddleware)

	makeDictionariesEndpoints(r, dictionaryService)
	makeConfigurationEndpoints(r, configurationService)
	makeCalendarEndpoints(r, calendarService)

	http.Handle("/", r)
	slog.Info("Started on port ", c.ServerPort)

	if c.EnableEureka {

		slog.Info("Registering into Eureka server")

		eurekaInstance, err := eureka.New().
			Default(c.ServerPort, "").
			Register(c.EurekaServiceUrl, "dictionary-service")

		util.FailOnError(err, "can't register with Eureka")

		defer eurekaInstance.Deregister()
	}

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

	db := pg.Connect(&pg.Options{
		User:       c.DatasourceUser,
		Password:   c.DatasourcePassword,
		Addr:       c.DatasourceHost,
		Database:   c.DatasourceName,
		PoolSize:   c.DatasourcePoolSize,
		MaxRetries: c.DatasourceMaxRetries,

		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			_, err := conn.Exec("set search_path=?", c.DatasourceSchema)
			if err != nil {
				slog.Error(err)
			}
			return nil
		},
	})

	for i := 0; i < c.InitDBConnectionRets; i++ {
		_, err := db.Exec("select 1")
		if err == nil {
			break
		}
		time.Sleep(1)
		slog.Info("trying to reconnect database")
	}

	if c.ShowSql {
		db.AddQueryHook(dbLogger{})
	}

	return db
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	slog.Debug(q.Query)
	return nil
}

func migrate(db *pg.DB) {
	flag.Parse()

	def := migrations.DefaultCollection
	def = def.SetTableName(fmt.Sprintf("%s.gopg_migrations", c.DatasourceSchema))

	if len(flag.Args()) == 0 {
		slog.Info("0 command line args ")
		_, _, err := def.Run(db, "init")
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
		ctx := tenant.NewContext(r.Context(), tenant.FromRequest(r))
		ctx = context.WithValue(ctx, "language", r.Header.Get("Accept-Language"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// access control and  CORS middleware
func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Tenant-Merge-Default, "+c.TenantHeaderName)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
