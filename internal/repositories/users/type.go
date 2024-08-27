package users

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"user-service/internal/model"
)

type FindOneUserInput struct {
	ID    null.Int
	Email null.String
}

type FindOneUserOutput struct {
	Data model.User
}

type CreateUserInput struct {
	Tx      wsqlx.Rdbms
	Payload model.User
}

type CreateUserOutput struct {
	ID int64
}

type CheckExistingUserInput struct {
	ID              null.Int
	Email           null.String
	IsEmailVerified null.Bool
}