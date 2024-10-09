package rabbitmq

import (
	"google.golang.org/protobuf/proto"
)

type PublishInput struct {
	RoutingKey string
	Exchange   string
	Payload    proto.Message
}
