package users

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"user-service/internal/repositories"
	"user-service/internal/util/tracer"
)

func (r *repository) FindOneUser(ctx context.Context, input FindOneUserInput) (output FindOneUserOutput, err error) {
	query := r.sq.Select("id", "is_email_verified", "email").From("users")
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}
	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeStruct, &output.Data)
	if err != nil {
		if errors.Is(err, wsqlx.ErrRecordNoRows) {
			err = repositories.ErrRecordNotFound
		}

		return output, tracer.Error(err)
	}

	return
}
