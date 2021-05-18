package main

import (
	"context"
	"dictionaries-service/service"
	"dictionaries-service/tenant"
	http2 "dictionaries-service/transport/http"
	"dictionaries-service/util"
	"flag"
	"fmt"
	"github.com/alapierre/gokit-utils/eureka"
	slog "github.com/go-eden/slf4go"
	httptransport "github.com/go-kit/kit/transport/http"
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
	ServerPort           int    `split_words:"true" default:"9098"`
	DatasourceName       string `split_words:"true" default:"app"`
	DatasourceSchema     string `split_words:"true" default:"dictionary"`
	DatasourceHost       string `split_words:"true" default:"localhost:5432"`
	DatasourceUser       string `required:"true" split_words:"true"`
	DatasourcePassword   string `required:"true" split_words:"true"`
	DatasourcePoolSize   int    `split_words:"true" default:"2"`
	DatasourceMaxRetries int    `split_words:"true" default:"3"`
	EurekaServiceUrl     string `split_words:"true" default:"http://localhost:8761/eureka"`
	DefaultLanguage      string `split_words:"true" default:"en"`
	InitDBConnectionRets int    `split_words:"true" default:"100"`
	ShowSql              bool   `split_words:"true" default:"true"`
	TenantHeaderName     string `split_words:"true" default:"X-Tenant-ID"`
}

var c Config

func main() {

	slog.Info("Starting up Dictionary Service")
	err := envconfig.Process("dict", &c)
	util.FailOnError(err, "Can't parse environment variables")

	slog.Infof("database name: %s host: %s user: %s", c.DatasourceName, c.DatasourceHost, c.DatasourceUser)

	http2.DefaultLanguage = language.MustParse(c.DefaultLanguage)
	slog.Info("Default language is ", http2.DefaultLanguage)

	tenant.HeaderName(c.TenantHeaderName)
	slog.Info("Tenant header is ", c.TenantHeaderName)

	db := connectDb()
	defer util.Close(db)

	migrate(db)

	dictionaryRepository := service.NewDictionaryRepository(db)
	translationRepository := service.NewTranslateRepository(db)
	metadataRepository := service.NewDictionaryMetadataRepository(db)
	dictionaryService := service.NewDictionaryService(dictionaryRepository, translationRepository, metadataRepository)

	configurationRepository := service.NewConfigurationRepository(db)
	configurationService := service.NewConfigurationService(configurationRepository)

	r := makeDictionariesEndpoints(dictionaryService)
	makeConfigurationEndpoints(r, configurationService)

	http.Handle("/", r)
	slog.Info("Started on port ", c.ServerPort)

	eurekaInstance, err := eureka.New().
		Default(c.ServerPort, "").
		Register(c.EurekaServiceUrl, "dictionary-service")

	util.FailOnError(err, "can't register with Eureka")

	defer eurekaInstance.Deregister()

	startHttpAndWaitForSigINT(c.ServerPort)

	slog.Info("Bye.")
}

func makeDictionariesEndpoints(dictionaryService *service.DictionaryService) *mux.Router {

	r := mux.NewRouter()
	r.Use(addContext, accessControlMiddleware)

	r.Methods("GET").Path("/api/dictionary/{type}/{key}").Handler(httptransport.NewServer(
		http2.MakeLoadDictEndpoint(dictionaryService),
		http2.DecodeLoadDictRequest,
		http2.EncodeResponse,
	))

	r.Methods("GET").Path("/api/metadata/{type}").Handler(httptransport.NewServer(
		http2.MakeLoadMetadataEndpoint(dictionaryService),
		http2.DecodeLoadMetadataRequest,
		http2.EncodeMetadataResponse,
	))

	r.Methods("GET").Path("/api/metadata").Handler(httptransport.NewServer(
		http2.MakeAvailableDictionaryTypesEndpoint(dictionaryService),
		func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		http2.EncodeResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/metadata").Handler(httptransport.NewServer(
		http2.MakeSaveMetadataEndpoint(dictionaryService),
		http2.DecodeSaveMetadataRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/metadata/{type}").Handler(httptransport.NewServer(
		http2.MakeSaveMetadataEndpointBetter(dictionaryService),
		http2.DecodeSaveMetadataRequestBetter,
		http2.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/metadata/{type}").Handler(httptransport.NewServer(
		http2.MakeUpdateMetadataEndpointBetter(dictionaryService),
		http2.DecodeSaveMetadataRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/metadata/{type}").Handler(httptransport.NewServer(
		http2.MakeSaveMetadataEndpointBetter(dictionaryService),
		http2.DecodeSaveMetadataRequestBetter,
		http2.EncodeSavedResponse,
	))

	r.Methods("GET").Path("/api/dictionary/{type}").Handler(httptransport.NewServer(
		http2.MakeLoadDictionaryByType(dictionaryService),
		http2.DecodeByTypeRequest,
		http2.EncodeResponse,
	))

	r.Methods("GET").Path("/api/dictionary/{type}/{key}/shallow").Handler(httptransport.NewServer(
		http2.MakeLoadDictShallowEndpoint(dictionaryService),
		http2.DecodeLoadDictRequest,
		http2.EncodeResponse,
	))

	r.Methods("GET").Path("/api/dictionary/{type}/{key}/children").Handler(httptransport.NewServer(
		http2.MakeLoadDictChildrenEndpoint(dictionaryService),
		http2.DecodeLoadDictRequest,
		http2.EncodeResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/dictionary").Handler(httptransport.NewServer(
		http2.MakeSaveDictionaryEndpoint(dictionaryService),
		http2.DecodeSaveDictRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/dictionary").Handler(httptransport.NewServer(
		http2.MakeUpdateDictionaryEndpoint(dictionaryService),
		http2.DecodeSaveDictRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("POST", "OPTIONS").Path("/api/dictionary/shallow").Handler(httptransport.NewServer(
		http2.MakeShallowSaveDictionaryEndpoint(dictionaryService),
		http2.DecodeShallowSaveDictionaryRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("PUT", "OPTIONS").Path("/api/dictionary/shallow").Handler(httptransport.NewServer(
		http2.MakeShallowUpdateDictionaryEndpoint(dictionaryService),
		http2.DecodeShallowSaveDictionaryRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("DELETE", "OPTIONS").Path("/api/dictionary/{type}/{key}").Handler(httptransport.NewServer(
		http2.MakeDeleteDictionaryEndpoint(dictionaryService),
		http2.DecodeLoadDictRequest,
		http2.EncodeSavedResponse,
	))

	r.Methods("DELETE").Path("/api/dictionary/all").Handler(httptransport.NewServer(
		http2.MakeDeleteAllDictionaryEndpoint(dictionaryService),
		func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		http2.EncodeSavedResponse,
	))

	r.Methods("DELETE").Path("/api/dictionary/{type}").Handler(httptransport.NewServer(
		http2.MakeDeleteDictionaryByTypeEndpoint(dictionaryService),
		http2.DecodeByTypeRequest,
		http2.EncodeSavedResponse,
	))
	return r
}

func makeConfigurationEndpoints(r *mux.Router, configurationService service.ConfigurationService) {

	r.Methods("GET").Path("/api/config/{key}/{day}").Handler(httptransport.NewServer(
		http2.MakeLoadConfigurationEndpoint(configurationService),
		http2.DecodeLoadConfigurationRequest,
		http2.EncodeResponse,
	))

	r.Methods("GET").Path("/api/configs/{day}").Handler(httptransport.NewServer(
		http2.MakeLoadConfigurationArrayEndpoint(configurationService),
		http2.DecodeLoadConfigurationArrayRequest,
		http2.EncodeResponse,
	))
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
