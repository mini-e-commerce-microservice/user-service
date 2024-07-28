package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"user-service/internal/conf"
	"user-service/internal/infra"
	"user-service/internal/presenter"
	"user-service/internal/provider/observability"
	"user-service/internal/util"
)

func main() {
	conf.Init()

	infra.OpenConnectionPostgres()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	closed := observability.InitOtel("0.0.0.0:4317", "user-svc")
	defer func() {
		err := closed(context.Background())
		util.Panic(err)
	}()

	presenter.New()
}
