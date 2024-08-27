package users

import "context"

type Repository interface {
	FindOneUser(ctx context.Context, input FindOneUserInput) (output FindOneUserOutput, err error)

	// CheckExistingUser if return true, data available
	CheckExistingUser(ctx context.Context, input CheckExistingUserInput) (exists bool, err error)
	CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
}
