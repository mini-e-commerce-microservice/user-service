package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories/users"
	jwt_util "github.com/mini-e-commerce-microservice/user-service/internal/util/jwt"
)

func (s *service) ActivationEmailUser(ctx context.Context, input ActivationEmailUserInput) (err error) {
	claims := jwt_util.OtpActivationEmailClaims{}

	err = claims.ClaimHS256(input.Token, s.jwtConf.OtpUsecase.Key)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = errors.Join(err, ErrTokenIsExpired)
		} else {
			err = fmt.Errorf("err from jwt: %s. %w", err.Error(), ErrInvalidToken)
		}
		return collection.Err(err)
	}

	err = s.userRepository.UpdateUser(ctx, users.UpdateUserInput{
		ID:    null.IntFrom(claims.UserId),
		Email: null.StringFrom(claims.Email),
		Payload: users.UpdateUserInputPayload{
			IsEmailVerified: null.BoolFrom(claims.IsVerified),
		},
	})
	if err != nil {
		return collection.Err(err)
	}
	return
}
