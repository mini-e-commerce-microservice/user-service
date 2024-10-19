package profiles

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
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
		return output, collection.Err(err)
	}
	return
}
