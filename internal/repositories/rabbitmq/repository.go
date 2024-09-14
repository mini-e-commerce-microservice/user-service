package rabbitmq

import (
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
)

type rabbitmq struct {
	client erabbitmq.RabbitMQPubSub
}

func NewRabbitMq(client erabbitmq.RabbitMQPubSub) *rabbitmq {
	return &rabbitmq{
		client: client,
	}
}
