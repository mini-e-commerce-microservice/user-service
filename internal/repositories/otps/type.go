package otps

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
)

type CreateOtpInput struct {
	Tx      wsqlx.Rdbms
	Payload model.Otp
}

type CreateOtpOutput struct {
	ID int64
}

type UpdateOtpInput struct {
	Tx      wsqlx.Rdbms
	ID      int64
	Payload UpdateOtpInputPayload
}

type UpdateOtpInputPayload struct {
	Token   null.String
	Counter null.Int16
}

type DeleteOtpInput struct {
	Tx         wsqlx.Rdbms
	ID         null.Int
	UserID     null.Int
	Usecase    null.String
	Type       null.String
	ExpiredGTE null.Time
	TokenIsNil bool
}

type FindOneOtpInput struct {
	ID         null.Int
	UserID     null.Int
	Code       null.String
	Usecase    null.String
	Type       null.String
	ExpiredGTE null.Time
	TokenIsNil bool
}

type FindOneOtpOutput struct {
	Data model.Otp
}
