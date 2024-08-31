package services

import (
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/jmoiron/sqlx"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/profiles"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/user-service/internal/services/user"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Dependency struct {
	UserService user.Service
}

func NewDependency(minioClient *minio.Client, db *sqlx.DB, rabbitmqClient *amqp.Channel) *Dependency {
	s3 := s3_wrapper_minio.New(minioClient)
	rdbms := wsqlx.NewRdbms(db)
	dbtx := wsqlx.NewSqlxTransaction(db)

	userRepository := users.NewRepository(rdbms)
	profileRepository := profiles.NewRepository(rdbms)
	rabbitmqRepository := rabbitmq.NewRabbitMq(rabbitmqClient)

	userSvc := user.NewService(user.NewServiceOptions{
		S3:                 s3,
		UserRepository:     userRepository,
		ProfileRepository:  profileRepository,
		RabbitmqRepository: rabbitmqRepository,
		DBTx:               dbtx,
	}, conf.GetConfig().Minio)

	return &Dependency{
		UserService: userSvc,
	}
}
