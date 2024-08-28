package rabbitmq_test

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
	"user-service/generated/proto/notification_proto"
	"user-service/internal/conf"
	"user-service/internal/infra"
	"user-service/internal/repositories/rabbitmq"
)

func TestRepo_Publish(t *testing.T) {
	t.Skip("integration test")
	conf.Init()
	_, ch, closed := infra.NewRabbitMq(conf.GetConfig().RabbitMQ)
	defer closed(context.Background())

	r := rabbitmq.NewRabbitMq(ch)
	err := r.Publish(context.TODO(), rabbitmq.PublishInput{
		RoutingKey: rabbitmq.RoutingKeyEmailOTP,
		Exchange:   rabbitmq.ExchangeNameNotification,
		Payload: &notification_proto.Notification{
			Type: notification_proto.NotificationType_EMAIL_VERIFIED,
			Data: &notification_proto.Notification_EmailVerified{
				EmailVerified: &notification_proto.NotificationEmailVerifiedPayload{
					OtpCode:   "123",
					ExpiredAt: timestamppb.New(time.Now().UTC()),
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	var forever chan struct{}

	msgs, err := ch.Consume("email_otp", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			var notification notification_proto.Notification
			err = proto.Unmarshal(d.Body, &notification)
			if err != nil {
				panic(err)
			}

			log.Printf("Received a message: %s", notification.Type)
			log.Printf("Received a correlation id: %s", d.CorrelationId)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			forever <- struct{}{}
			break
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
