package errors

import "errors"

var (
	ErrorWrongPassword = errors.New("wrong password")
	ErrorWrongOTP      = errors.New("wrong one-time password")
)
