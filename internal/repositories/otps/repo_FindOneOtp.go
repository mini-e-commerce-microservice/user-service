package otps

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories"
)

func (r *repository) FindOneOtp(ctx context.Context, input FindOneOtpInput) (output FindOneOtpOutput, err error) {
	query := r.sq.Select(
		"id", "user_id", "usecase", "code", "type", "counter", "expired", "token",
	).From("otps").Limit(1)

	if input.Type.Valid {
		query = query.Where(squirrel.Eq{"type": input.Type.String})
	}

	if input.UserID.Valid {
		query = query.Where(squirrel.Eq{"user_id": input.UserID.Int64})
	}

	if input.Usecase.Valid {
		query = query.Where(squirrel.Eq{"usecase": input.Usecase.String})
	}

	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	if input.ExpiredGTE.Valid {
		query = query.Where(squirrel.GtOrEq{"expired": input.ExpiredGTE})
	}

	if input.Code.Valid {
		query = query.Where(squirrel.Eq{"code": input.Code.String})
	}

	if input.TokenIsNil {
		query = query.Where(squirrel.Eq{"token": nil})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeStruct, &output.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = repositories.ErrRecordNotFound
		}

		return output, collection.Err(err)
	}

	return
}
