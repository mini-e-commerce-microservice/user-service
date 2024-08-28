package users

import (
	"context"
	"github.com/Masterminds/squirrel"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"go.opentelemetry.io/otel"
	"user-service/internal/util/tracer"
)

func (r *repository) CheckExistingUser(ctx context.Context, input CheckExistingUserInput) (exists bool, err error) {
	ctx, span := otel.Tracer("users repository").Start(ctx, "check existing user data")
	defer span.End()

	query := r.sq.Select("1").
		Prefix("SELECT EXISTS(").
		Suffix(")").
		From("users")

	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}
	if input.IsEmailVerified.Valid {
		query = query.Where(squirrel.Eq{"is_email_verified": input.IsEmailVerified.Bool})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &exists)
	if err != nil {
		return exists, tracer.Error(err)
	}
	return
}
