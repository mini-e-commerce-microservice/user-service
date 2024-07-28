package infra

import (
	"context"
	pgxotel "github.com/SyaibanAhmadRamadhan/pgx-otel"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-service/internal/conf"
	"user-service/internal/util"
)

func OpenConnectionPostgres() *pgxpool.Pool {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conf.GetDatabaseDSN())
	util.Panic(err)

	pool.Config().ConnConfig.Tracer = pgxotel.NewTracer()

	err = pool.Ping(ctx)
	util.Panic(err)

	return pool
}
