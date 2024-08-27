package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"os/signal"
	"syscall"
	"user-service/internal/conf"
	"user-service/internal/infra"
	"user-service/internal/presenter"
)

var restApiCmd = &cobra.Command{
	Use:   "rest-api",
	Short: "run rest api",
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init()

		otel := infra.NewOtel(conf.GetConfig().OpenTelemetry)
		_, dbClose := infra.NewPostgresql(conf.GetConfig().DatabaseDSN)

		server := presenter.New(&presenter.Presenter{
			Port: conf.GetConfig().AppPort,
		})

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := server.Shutdown(context.TODO()); err != nil {
			panic(err)
		}

		if err := otel(context.TODO()); err != nil {
			panic(err)
		}

		if err := dbClose(context.TODO()); err != nil {
			panic(err)
		}

		log.Info().Msg("Shutdown complete. Exiting.")
		return
	},
}
