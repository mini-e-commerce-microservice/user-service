package otps

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
)

func (r *repository) UpdateOtp(ctx context.Context, input UpdateOtpInput) (err error) {
	query := r.sq.Update("otps").Where(squirrel.Eq{"id": input.ID})
	if input.Payload.Token.Valid {
		query = query.Set("token", input.Payload.Token.String)
	}

	if input.Payload.Counter.Valid {
		query = query.Set("counter", input.Payload.Counter.Int16)
	}

	rdbms := r.rdbms
	if input.Tx != nil {
		rdbms = input.Tx
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return tracer.Error(err)
	}

	return
}
