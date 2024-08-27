package model

import (
	"time"
)

type User struct {
	ID              int64      `db:"id"`
	Email           string     `db:"email"`
	Password        string     `db:"password"`
	CreatedAt       time.Time  `db:"created_at"`
	IsEmailVerified bool       `db:"is_email_verified"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}
