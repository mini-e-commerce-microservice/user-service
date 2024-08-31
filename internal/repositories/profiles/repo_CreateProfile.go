package profiles

import (
	"context"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/user-service/internal/util/tracer"
	"time"
)

func (r *repository) CreateProfile(ctx context.Context, input CreateProfileInput) (output CreateProfileOutput, err error) {
	query := r.sq.Insert("profiles").Columns(
		"user_id", "full_name", "image_profile", "background_image", "created_at", "updated_at",
	).Values(
		input.Payload.UserID,
		input.Payload.FullName,
		input.Payload.ImageProfile,
		input.Payload.BackgroundImage,
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
