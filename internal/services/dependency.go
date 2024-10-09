package services

import (
	"context"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/infra"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/otp"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"github.com/rs/zerolog/log"
)

type Dependency struct {
	UserService user.Service
	OtpService  otp.Service
}

func NewDependency(appConf *conf.AppConfig) (*Dependency, func()) {
	otelConf := conf.LoadOtelConf()
	jwtConf := conf.LoadJwtConf()
	minioConf := conf.LoadMinioConf()
	rabbitMQConf := conf.LoadRabbitMQConf()

	otel := infra.NewOtel(otelConf, appConf.ServiceName)
	postgre, dbClose := infra.NewPostgresql(appConf.DatabaseDSN)
	minio := infra.NewMinio(conf.LoadMinioConf())
	rabbitmqClient := erabbitmq.New(rabbitMQConf.Url, erabbitmq.WithOtel(rabbitMQConf.Url))

	s3 := s3_wrapper_minio.New(minio)
	rdbms := wsqlx.NewRdbms(postgre)

	userRepository := users.NewRepository(rdbms)
	profileRepository := profiles.NewRepository(rdbms)
	otpRepository := otps.NewRepository(rdbms)
	rabbitmqRepository := rabbitmq.NewRabbitMq(rabbitmqClient)

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

	closeFn := func() {
		if err := dbClose(context.Background()); err != nil {
			log.Err(err).Msg("failed closed db")
		}

		rabbitmqClient.Close()

		if err := otel(context.Background()); err != nil {
			log.Err(err).Msg("failed closed otel")
		}
	}
	return &Dependency{
		UserService: userSvc,
		OtpService:  otpSvc,
	}, closeFn
}
