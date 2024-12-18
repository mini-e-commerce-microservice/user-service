package main

import (
	"context"
	"errors"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/presentations"
	"github.com/mini-e-commerce-microservice/user-service/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"os/signal"
	"syscall"
)

var restApiCmd = &cobra.Command{
	Use:   "restApi",
	Short: "run rest api",
	Run: func(cmd *cobra.Command, args []string) {
		appConf := conf.LoadAppConf()
		jwtConf := conf.LoadJwtConf()

		dependency, closeFn := services.NewDependency(appConf, jwtConf)

		server := presentations.New(&presentations.Presenter{
			Dependency:         dependency,
			JwtAccessTokenConf: jwtConf.AccessToken,
			Port:               int(appConf.AppPort),
		})

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Err(err).Msg("failed listen serve")
				stop()
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("failed shutdown server")
		}

		closeFn(context.Background())
		log.Info().Msg("Shutdown complete. Exiting.")
		return
	},
}
