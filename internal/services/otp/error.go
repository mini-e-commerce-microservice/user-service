package otp

import "errors"

var ErrDestinationAddressNotFound = errors.New("destination address not found")
var ErrOtpNotFound = errors.New("otp not found")
var ErrOtpExpired = errors.New("otp expired")
var ErrOtpCounterExceeded = errors.New("otp counter exceeded")
var ErrCodeOtpInvalid = errors.New("invalid code otp")
var ErrEmailUserIsVerified = errors.New("email user has been verified")
