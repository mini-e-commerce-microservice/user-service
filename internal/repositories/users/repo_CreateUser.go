package users

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
	"time"
)

func (r *repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	query := r.sq.Insert("users").Columns(
		"email", "password", "is_email_verified", "created_at", "updated_at", "trace_parent",
	).Values(
		input.Payload.Email,
		input.Payload.Password,
		input.Payload.IsEmailVerified,
		time.Now().UTC(),
		time.Now().UTC(),
		util.GetTraceParent(ctx),
	).Suffix("RETURNING id")

	rdbms := r.rdbms
	if input.Tx != nil {
		rdbms = input.Tx
	}

	err = rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &output.ID)
	if err != nil {
		return output, collection.Err(err)
	}

	return
}
