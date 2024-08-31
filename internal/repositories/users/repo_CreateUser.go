package users

import (
	"context"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"time"
)

func (r *repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	query := r.sq.Insert("users").Columns(
		"email", "password", "is_email_verified", "register_as", "created_at", "updated_at",
	).Values(
		input.Payload.Email,
		input.Payload.Password,
		input.Payload.IsEmailVerified,
		input.Payload.RegisterAs,
		time.Now().UTC(),
		time.Now().UTC(),
	).Suffix("RETURNING id")

	rdbms := r.rdbms
	if input.Tx != nil {
		rdbms = input.Tx
	}

	err = rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &output.ID)
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}
