package profiles

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/user-service/internal/model"
	"github.com/mini-e-commerce-microservice/user-service/internal/repositories"
)

func (r *repository) FindOneProfile(ctx context.Context, input FindOneProfileInput) (output FindOneProfileOutput, err error) {
	query := r.sq.Select("id", "user_id", "full_name", "image_profile", "background_image").
		From("profiles")
	if input.UserID.Valid {
		query = query.Where(squirrel.Eq{"user_id": input.UserID.Int64})
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

type FindOneProfileInput struct {
	UserID null.Int
}
type FindOneProfileOutput struct {
	Data model.Profile
}
