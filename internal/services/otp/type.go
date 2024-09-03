package otp

import "github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"

type SendOtpInput struct {
	Usecase            primitive.OtpUseCase
	Type               primitive.OtpType
	DestinationAddress string
}

type VerifyOtpInput struct {
	Usecase            primitive.OtpUseCase
	Type               primitive.OtpType
	DestinationAddress string
	Code               string
}

type VerifyOtpOutput struct {
	Token string
}
