package util

import "github.com/go-eden/slf4go"

func FailOnError(err error, msg string) {
	if err != nil {
		slog.Fatalf("%s: %s", msg, err)
		panic(err)
	}
}

type Closable interface {
	Close() error
}

func Close(target Closable) {
	err := target.Close()
	if err != nil {
		slog.Errorf("Cant close %v", target)
	}
}
