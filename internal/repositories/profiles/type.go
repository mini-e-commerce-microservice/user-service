package profiles

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
)

type CreateProfileInput struct {
	Tx      wsqlx.Rdbms
	Payload model.Profile
}

type CreateProfileOutput struct {
	ID int64
}
