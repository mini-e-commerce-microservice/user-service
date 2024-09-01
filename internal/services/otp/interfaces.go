package otp

import "context"

type Service interface {
	SendOtp(ctx context.Context, input SendOtpInput) (err error)
	VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error)
}
