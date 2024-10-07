package jwt_util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/user-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"time"
)

type OtpActivationEmailClaims struct {
	*jwt_claims_proto.JwtOtpActivationEmailClaims
	jwt.RegisteredClaims
}

func (o *OtpActivationEmailClaims) GenerateHS256(key string, exp time.Duration) (string, error) {
	timeNow := time.Now().UTC()
	timeExp := timeNow.Add(exp)

	o.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(timeExp),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		Subject:   "activation_email",
		Audience:  []string{"activation_email"},
	}
	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, o)

	tokenStr, err := tokenParse.SignedString([]byte(key))
	if err != nil {
		return "", tracer.Error(err)
	}

	return tokenStr, nil
}

func (o *OtpActivationEmailClaims) ClaimHS256(tokenStr, key string) error {
	tokenParse, err := jwt.ParseWithClaims(tokenStr, o, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return tracer.Error(err)
	}

	if !tokenParse.Valid {
		return tracer.Error(ErrInvalidParseToken)
	}

	return nil
}
