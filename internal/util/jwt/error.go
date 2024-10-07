package jwt_util

import (
	"errors"
)

var ErrInvalidClaims = errors.New("invalid claims")
var ErrInvalidParseToken = errors.New("invalid parse token")
