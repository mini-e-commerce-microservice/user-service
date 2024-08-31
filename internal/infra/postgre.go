package infra

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mini-e-commerce-microservice/user-service/internal/primitive"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"github.com/rs/zerolog/log"
	"time"
)

func NewPostgresql(dsn string) (*sqlx.DB, primitive.CloseFn) {
	db, err := sqlx.Connect("postgres", dsn)
	util.Panic(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	util.Panic(err)

	log.Info().Msg("initialization postgresql successfully")
	return db, func(ctx context.Context) (err error) {
		log.Info().Msg("starting close postgresql db")

		err = db.Close()
		if err != nil {
			return err
		}

		log.Info().Msg("close postgresql db successfully")
		return
	}
}
