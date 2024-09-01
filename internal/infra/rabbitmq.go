package infra

import (
	"context"
	"github.com/mini-e-commerce-microservice/user-service/internal/conf"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func NewRabbitMq(cred conf.ConfigRabbitMQ) (*amqp.Connection, *amqp.Channel, primitive.CloseFn) {
	conn, err := amqp.Dial(cred.Url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	log.Info().Msg("initialization rabbitmq successfully")
	return conn, ch, func(ctx context.Context) (err error) {
		err = conn.Close()
		if err != nil {
			log.Err(err).Msg("failed close rabbitmq connection")
		}

		log.Info().Msg("closed rabbitmq connection successfully")

		err = ch.Close()
		if err != nil {
			log.Err(err).Msg("failed close rabbitmq channel")
		}

		log.Info().Msg("closed rabbitmq channel successfully")
		return nil
	}

}
