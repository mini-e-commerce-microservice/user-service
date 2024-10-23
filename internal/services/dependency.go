package services

import (
	"context"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/infra"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/otp"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"github.com/rs/zerolog/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Dependency struct {
	UserService user.Service
	OtpService  otp.Service
}

func NewDependency(appConf *secret_proto.UserService, jwtConf *secret_proto.Jwt) (*Dependency, func(ctx context.Context)) {
	otelConf := conf.LoadOtelConf()
	minioConf := conf.LoadMinioConf()
	rabbitMQConf := conf.LoadRabbitMQConf()

	otelCloseFn := infra.NewOtel(otelConf, appConf.TracerName)
	pgdb, pgdbCloseFn := infra.NewPostgresql(appConf.DatabaseDsn)
	minio := infra.NewMinio(conf.LoadMinioConf())
	rabbitmqClient := erabbitmq.New(rabbitMQConf.Url, erabbitmq.WithOtel(rabbitMQConf.Url))

	s3 := s3_wrapper_minio.New(minio)
	rdbms := wsqlx.NewRdbms(pgdb, wsqlx.WithAttributes(semconv.DBSystemPostgreSQL))

	// REPOSITORY LAYER
	userRepository := users.NewRepository(rdbms)
	profileRepository := profiles.NewRepository(rdbms)
	otpRepository := otps.NewRepository(rdbms)
	rabbitmqRepository := rabbitmq.NewRabbitMq(rabbitmqClient)

	// SERVICE LAYER
	userSvc := user.NewService(user.NewServiceOptions{
		S3:                 s3,
		UserRepository:     userRepository,
		ProfileRepository:  profileRepository,
		RabbitmqRepository: rabbitmqRepository,
		DBTx:               rdbms,
		BucketName:         minioConf.PrivateBucket,
		JwtConfig:          jwtConf,
	})
	otpSvc := otp.NewService(otp.NewServiceOptions{
		UserRepository:     userRepository,
		RabbitmqRepository: rabbitmqRepository,
		OtpRepository:      otpRepository,
		DBTx:               rdbms,
		JwtConf:            jwtConf,
		RabbitMQConf:       rabbitMQConf,
	})

	closeFn := func(ctx context.Context) {
		if err := pgdbCloseFn(ctx); err != nil {
			log.Err(err).Msg("failed closed db")
		}

		rabbitmqClient.Close()

		if err := otelCloseFn(ctx); err != nil {
			log.Err(err).Msg("failed closed otel")
		}
	}
	return &Dependency{
		UserService: userSvc,
		OtpService:  otpSvc,
	}, closeFn
}
