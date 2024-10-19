package otps

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
)

func (r *repository) DeleteOtp(ctx context.Context, input DeleteOtpInput) (err error) {
	query := r.sq.Delete("otps")

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

	if input.TokenIsNil {
		query = query.Where(squirrel.Eq{"token": nil})
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
