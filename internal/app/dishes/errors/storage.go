package errors

import "errors"

var (
	ErrorUniqueViolationWhileStoring = errors.New("ingredient with this name already exists")
	ErrorNoRowsWhileGetting          = errors.New("ingredient does not exist")
	ErrorNoRowsUpdated               = errors.New("ingredient does not exist")
)
