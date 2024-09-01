package model

import "time"

type Otp struct {
	ID      int64     `db:"id"`
	UserID  int64     `db:"user_id"`
	Usecase string    `db:"usecase"`
	Code    string    `db:"code"`
	Type    string    `db:"type"`
	Counter int16     `db:"counter"`
	Expired time.Time `db:"expired"`
	Token   *string   `db:"token"`
}
