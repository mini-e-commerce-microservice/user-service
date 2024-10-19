package users

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/user-service/internal/util"
)

func (r *repository) UpdateUser(ctx context.Context, input UpdateUserInput) (err error) {
	query := r.sq.Update("users").Set("trace_parent", util.GetTraceParent(ctx))
	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	if input.Payload.IsEmailVerified.Valid {
		query = query.Set("is_email_verified", input.Payload.IsEmailVerified.Bool)
	}
	if input.Payload.Password.Valid {
		query = query.Set("password", input.Payload.Password.String)
	}

	rdbms := r.rdbms
	if input.Tx != nil {
		rdbms = input.Tx
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}
	return
}
