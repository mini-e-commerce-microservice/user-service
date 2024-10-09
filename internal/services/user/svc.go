package user

import (
	s3wrapper "github.com/SyaibanAhmadRamadhan/go-s3-wrapper"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
)

type service struct {
	s3                 s3wrapper.S3Client
	userRepository     users.Repository
	profileRepository  profiles.Repository
	rabbitmqRepository rabbitmq.Repository
	dbTx               wsqlx.Tx
	bucketName         string
	jwtConf            *secret_proto.Jwt
}

type NewServiceOptions struct {
	S3                 s3wrapper.S3Client
	UserRepository     users.Repository
	ProfileRepository  profiles.Repository
	RabbitmqRepository rabbitmq.Repository
	DBTx               wsqlx.Tx
	BucketName         string
	JwtConfig          *secret_proto.Jwt
}

func NewService(opts NewServiceOptions) *service {
	return &service{
		s3:                 opts.S3,
		userRepository:     opts.UserRepository,
		profileRepository:  opts.ProfileRepository,
		rabbitmqRepository: opts.RabbitmqRepository,
		dbTx:               opts.DBTx,
		bucketName:         opts.BucketName,
		jwtConf:            opts.JwtConfig,
	}
}
