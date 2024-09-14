package rabbitmq

import (
	"context"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	httplogwrap "github.com/SyaibanAhmadRamadhan/http-log-wrap"
	"github.com/google/uuid"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func (r *rabbitmq) Publish(ctx context.Context, input PublishInput) (err error) {
	correlationID := httplogwrap.GetCorrelationID(ctx)
	messageID := uuid.New().String()

	body, err := proto.Marshal(input.Payload)
	if err != nil {
		return tracer.Error(err)
	}

	_, err = r.client.Publish(
		ctx,
		erabbitmq.PubInput{
			ExchangeName: string(input.Exchange),
			RoutingKey:   string(input.RoutingKey),
			Mandatory:    false,
			Immediate:    false,
			Msg: amqp.Publishing{
				MessageId:     messageID,
				CorrelationId: correlationID,
				ContentType:   "application/protobuf",
				Body:          body,
				Headers: amqp.Table{
					"correlation_id": correlationID,
				},
			},
			DelayRetry: 0,
			MaxRetry:   1,
		},
	)
	if err != nil {
		return tracer.Error(err)
	}

	return
}
