package rabbitmq

import (
	"google.golang.org/protobuf/proto"
)

type RoutingKey string

const (
	RoutingKeyNotificationTypeEmail RoutingKey = "notification.type.email"
)

type Exchange string

const (
	ExchangeNameNotification Exchange = "notifications"
)

type PublishInput struct {
	RoutingKey RoutingKey
	Exchange   Exchange
	Payload    proto.Message
}
