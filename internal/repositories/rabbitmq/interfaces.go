package rabbitmq

import "context"

type Repository interface {
	Publish(ctx context.Context, input PublishInput) (err error)
}
