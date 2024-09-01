package otp

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
)

type service struct {
	userRepository     users.Repository
	rabbitmqRepository rabbitmq.Repository
	otpRepository      otps.Repository
	dbTx               wsqlx.Tx
	jwtKey             string
}

type NewServiceOptions struct {
	UserRepository     users.Repository
	RabbitmqRepository rabbitmq.Repository
	OtpRepository      otps.Repository
	DBTx               wsqlx.Tx
}

func NewService(opts NewServiceOptions, cred conf.ConfigJWT) *service {
	return &service{
		userRepository:     opts.UserRepository,
		otpRepository:      opts.OtpRepository,
		rabbitmqRepository: opts.RabbitmqRepository,
		dbTx:               opts.DBTx,
		jwtKey:             cred.OtpUsecase,
	}
}
