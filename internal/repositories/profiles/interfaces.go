package profiles

import "context"

type Repository interface {
	CreateProfile(ctx context.Context, input CreateProfileInput) (output CreateProfileOutput, err error)
	FindOneProfile(ctx context.Context, input FindOneProfileInput) (output FindOneProfileOutput, err error)
}
