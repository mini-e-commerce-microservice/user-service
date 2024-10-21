package user

import "context"

type Service interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (output RegisterUserOutput, err error)
	ActivationEmailUser(ctx context.Context, input ActivationEmailUserInput) (err error)
	GetUser(ctx context.Context, input GetUserInput) (output GetUserOutput, err error)
}
