package users

import (
	"context"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"time"
	"user-service/internal/util/tracer"
)

func (r *repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	query := r.sq.Insert("users").Columns(
		"email", "password", "is_email_verified", "created_at", "updated_at",
	).Values(
		input.Payload.Email,
		input.Payload.Password,
		input.Payload.IsEmailVerified,
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
