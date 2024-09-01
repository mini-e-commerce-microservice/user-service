package otps

import "context"

type Repository interface {
	CreateOtp(ctx context.Context, input CreateOtpInput) (output CreateOtpOutput, err error)
	UpdateOtp(ctx context.Context, input UpdateOtpInput) (err error)
	DeleteOtp(ctx context.Context, input DeleteOtpInput) (err error)
	FindOneOtp(ctx context.Context, input FindOneOtpInput) (output FindOneOtpOutput, err error)
}
