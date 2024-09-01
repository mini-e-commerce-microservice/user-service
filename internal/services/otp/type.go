package otp

import "github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"

type SendOtpInput struct {
	Usecase            primitive.OtpUseCase
	UserID             int64
	Type               primitive.OtpType
	DestinationAddress string
}

type VerifyOtpInput struct {
	Usecase primitive.OtpUseCase
	Type    primitive.OtpType
	Code    string
	UserID  int64
}

type VerifyOtpOutput struct {
	Token string
}
