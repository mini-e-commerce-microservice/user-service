package services

import (
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/jmoiron/sqlx"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/otp"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"github.com/minio/minio-go/v7"
)

type Dependency struct {
	UserService user.Service
	OtpService  otp.Service
}

func NewDependency(minioClient *minio.Client, db *sqlx.DB, rabbitmqClient erabbitmq.RabbitMQPubSub) *Dependency {
	s3 := s3_wrapper_minio.New(minioClient)
	rdbms := wsqlx.NewRdbms(db)

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
		MinioConfig:        conf.GetConfig().Minio,
		JwtConfig:          conf.GetConfig().Jwt,
	})
	otpSvc := otp.NewService(otp.NewServiceOptions{
		UserRepository:     userRepository,
		RabbitmqRepository: rabbitmqRepository,
		OtpRepository:      otpRepository,
		DBTx:               rdbms,
	}, conf.GetConfig().Jwt)

	return &Dependency{
		UserService: userSvc,
		OtpService:  otpSvc,
	}
}
