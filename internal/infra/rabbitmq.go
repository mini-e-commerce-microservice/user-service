package infra

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"user-service/internal/conf"
	"user-service/internal/primitive"
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
