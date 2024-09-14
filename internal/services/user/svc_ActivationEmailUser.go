package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/jwt_payload_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	jwt_util "github.com/mini-e-commerce-microservice/user-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"google.golang.org/protobuf/proto"
)

func (s *service) ActivationEmailUser(ctx context.Context, input ActivationEmailUserInput) (err error) {
	claims, err := jwt_util.ClaimHS256(input.Token, s.keyJwtOtp)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = errors.Join(err, ErrTokenIsExpired)
		} else {
			err = fmt.Errorf("err from jwt: %s. %w", err.Error(), ErrInvalidToken)
		}
		return tracer.Error(err)
	}

	payloadStr, ok := claims["payload"].(string)
	if !ok {
		return tracer.Error(ErrInvalidToken)
	}

	payload := &jwt_payload_proto.JwtPayload{}

	err = proto.Unmarshal([]byte(payloadStr), payload)
	if err != nil {
		err = errors.Join(err, ErrInvalidToken)
		return tracer.Error(err)
	}

	verifyEmail, ok := payload.Data.(*jwt_payload_proto.JwtPayload_OtpActivationEmail)
	if !ok {
		return tracer.Error(ErrInvalidToken)
	}

	err = s.userRepository.UpdateUser(ctx, users.UpdateUserInput{
		ID:    null.IntFrom(verifyEmail.OtpActivationEmail.UserId),
		Email: null.StringFrom(verifyEmail.OtpActivationEmail.Email),
		Payload: users.UpdateUserInputPayload{
			IsEmailVerified: null.BoolFrom(verifyEmail.OtpActivationEmail.IsVerified),
		},
	})
	if err != nil {
		return tracer.Error(err)
	}
	return
}
