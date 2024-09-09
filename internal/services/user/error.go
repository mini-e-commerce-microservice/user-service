package user

import "errors"

var ErrEmailAvailable = errors.New("email is available")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenIsExpired = errors.New("token is expired")
