package test

import (
	"context"
	"github.com/go-eden/slf4go"
	"github.com/go-pg/pg/v10"
)

func ConnectDb() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "app",
		Password: "qwedsazxc",
		Addr:     "localhost:5432",
		Database: "app",
		OnConnect: func(ctx context.Context, conn *pg.Conn) error {
			_, err := conn.Exec("set search_path=?", "dictionary")
			if err != nil {
				slog.Error(err)
			}
			return nil
		},
	})
	db.AddQueryHook(dbLogger{})
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
