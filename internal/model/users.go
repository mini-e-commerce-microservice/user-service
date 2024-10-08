package model

import (
	"context"
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
	"time"
)

type User struct {
	ID              int64      `db:"id"`
	Email           string     `db:"email"`
	Password        string     `db:"password"`
	CreatedAt       time.Time  `db:"created_at"`
	IsEmailVerified bool       `db:"is_email_verified"`
	RegisterAs      int8       `db:"register_as"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

func GetTraceParent(ctx context.Context) *string {
	traceParent := whttp.GetTraceParent(ctx)
	if traceParent != "" {
		return &traceParent
	}

	return nil
}
