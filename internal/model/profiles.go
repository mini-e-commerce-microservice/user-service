package model

import (
	"time"
)

type Profile struct {
	ID              int64      `db:"id"`
	UserID          int64      `db:"user_id"`
	FullName        string     `db:"full_name"`
	ImageProfile    *string    `db:"image_profile"`
	BackgroundImage *string    `db:"background_image"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}
