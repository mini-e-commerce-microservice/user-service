package otp

import (
	"context"
	"errors"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/otp_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/otps"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	jwt_util "github.com/mini-e-commerce-microservice/user-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"google.golang.org/protobuf/proto"
	"time"
)

func (s *service) VerifyOtp(ctx context.Context, input VerifyOtpInput) (output VerifyOtpOutput, err error) {
	userOutput, err := s.validateExistingUser(ctx, input.Type, input.DestinationAddress)
	if err != nil {
		return output, tracer.Error(err)
	}

	otpOutput, err := s.otpRepository.FindOneOtp(ctx, otps.FindOneOtpInput{
		UserID:     null.IntFrom(userOutput.Data.ID),
		Usecase:    null.StringFrom(string(input.Usecase)),
		Type:       null.StringFrom(string(input.Type)),
		TokenIsNil: true,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			err = errors.Join(err, ErrOtpNotFound)
		}
		return output, tracer.Error(err)
	}

	if otpOutput.Data.Expired.Before(time.Now().UTC()) {
		return output, tracer.Error(ErrOtpExpired)
	}
	if otpOutput.Data.Counter >= input.Usecase.GetLimitRetry() {
		return output, tracer.Error(ErrOtpCounterExceeded)
	}

	token := null.String{}
	counter := null.Int16{}
	var errOtp error
	tokenStr := ""

	if otpOutput.Data.Code != input.Code {
		errOtp = ErrCodeOtpInvalid
		counter = null.Int16From(otpOutput.Data.Counter + 1)
	} else {
		tokenStr, err = s.generateTokenOTP(input, userOutput)
		if err != nil {
			return output, tracer.Error(err)
		}

		token = null.StringFrom(tokenStr)
	}

	err = s.otpRepository.UpdateOtp(ctx, otps.UpdateOtpInput{
		ID: otpOutput.Data.ID,
		Payload: otps.UpdateOtpInputPayload{
			Token:   token,
			Counter: counter,
		},
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	if errOtp != nil {
		return output, tracer.Error(errOtp)
	}

	output = VerifyOtpOutput{
		Token: tokenStr,
	}
	return
}

func (s *service) generateTokenOTP(input VerifyOtpInput, user users.FindOneUserOutput) (string, error) {
	payloadProto := &otp_proto.OTP{}

	if input.Usecase == primitive.OtpUseCaseVerifyEmail {
		payloadProto = &otp_proto.OTP{
			Payload: &otp_proto.OTP_UserVerifyEmail{
				UserVerifyEmail: &otp_proto.OtpUserVerifyEmailPayload{
					UserId:     user.Data.ID,
					Email:      user.Data.Email,
					IsVerified: true,
				},
			},
		}
	}

	body, err := proto.Marshal(payloadProto)
	if err != nil {
		return "", tracer.Error(err)
	}

	tokenStr, err := jwt_util.GenerateHS256(jwt_util.Jwt{
		UserID:  user.Data.ID,
		Key:     s.jwtKey,
		Payload: string(body),
		Exp:     input.Usecase.GetTTL(),
	})
	if err != nil {
		return "", tracer.Error(err)
	}

	return tokenStr, nil
}
