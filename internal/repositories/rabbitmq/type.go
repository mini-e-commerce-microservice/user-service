package rabbitmq

import (
	"google.golang.org/protobuf/proto"
)

type RoutingKey string

const (
	RoutingKeyEmailOTP RoutingKey = "notification.email.otp"
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
