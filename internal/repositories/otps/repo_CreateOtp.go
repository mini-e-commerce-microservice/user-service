package otps

import (
	"context"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
)

func (r *repository) CreateOtp(ctx context.Context, input CreateOtpInput) (output CreateOtpOutput, err error) {
	query := r.sq.Insert("otps").Columns(
		"user_id", "usecase", "code", "type", "counter", "expired",
	).Values(
		input.Payload.UserID,
		input.Payload.Usecase,
		input.Payload.Code,
		input.Payload.Type,
		input.Payload.Counter,
		input.Payload.Expired,
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
