package errors

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorNoSuchUser        = errors.New("no such user")
	ErrorWrongPassword     = errors.New("wrong password")
)
