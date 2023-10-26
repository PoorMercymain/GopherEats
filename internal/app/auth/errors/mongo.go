package errors

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorNoSuchUser        = errors.New("no such user")
	ErrorUserWasNotUpdated = errors.New("user was not updated")
)
