package jwt_util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"time"
)

type Jwt struct {
	UserID  int64
	Key     string
	Exp     time.Duration
	Payload string
}

func GenerateHS256(jwtModel Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     timeExp,
		"sub":     jwtModel.UserID,
		"payload": jwtModel.Payload,
	})

	tokenStr, err := tokenParse.SignedString([]byte(jwtModel.Key))
	if err != nil {
		return "", tracer.Error(err)
	}

	return tokenStr, nil
}
