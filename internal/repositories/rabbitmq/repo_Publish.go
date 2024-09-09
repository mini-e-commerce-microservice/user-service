package rabbitmq

import (
	"context"
	httplogwrap "github.com/SyaibanAhmadRamadhan/http-log-wrap"
	"github.com/google/uuid"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

func (r *rabbitmq) Publish(ctx context.Context, input PublishInput) (err error) {
	correlationID := httplogwrap.GetCorrelationID(ctx)
	messageID := uuid.New().String()

	ctx, span := otel.Tracer("rabbitmq").Start(ctx, "publish message", trace.WithAttributes(
		attribute.String("rabbitmq.message_id", messageID),
		attribute.String("rabbitmq.correlation_id", correlationID),
		attribute.String("rabbitmq.exchange", string(input.Exchange)),
		attribute.String("rabbitmq.routing_key", string(input.RoutingKey)),
	))
	defer span.End()

	body, err := proto.Marshal(input.Payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return tracer.Error(err)
	}

	err = r.ch.PublishWithContext(
		ctx,
		string(input.Exchange),
		string(input.RoutingKey),
		true,
		false,
		amqp.Publishing{
			MessageId:     messageID,
			CorrelationId: correlationID,
			ContentType:   "application/protobuf",
			Body:          body,
		},
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return tracer.Error(err)
	}

	return
}
