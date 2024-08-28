package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmq struct {
	ch *amqp.Channel
}

func NewRabbitMq(ch *amqp.Channel) *rabbitmq {
	return &rabbitmq{
		ch: ch,
	}
}
