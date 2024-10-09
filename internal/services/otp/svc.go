package otp

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
)

type service struct {
	userRepository     users.Repository
	rabbitmqRepository rabbitmq.Repository
	otpRepository      otps.Repository
	dbTx               wsqlx.Tx
	jwtConf            *secret_proto.Jwt
	rabbitMQConf       *secret_proto.RabbitMQ
}

type NewServiceOptions struct {
	UserRepository     users.Repository
	RabbitmqRepository rabbitmq.Repository
	OtpRepository      otps.Repository
	DBTx               wsqlx.Tx
	JwtConf            *secret_proto.Jwt
	RabbitMQConf       *secret_proto.RabbitMQ
}

func NewService(opts NewServiceOptions) *service {
	return &service{
		userRepository:     opts.UserRepository,
		otpRepository:      opts.OtpRepository,
		rabbitmqRepository: opts.RabbitmqRepository,
		dbTx:               opts.DBTx,
		jwtConf:            opts.JwtConf,
		rabbitMQConf:       opts.RabbitMQConf,
	}
}
