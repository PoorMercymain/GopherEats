package errors

import "errors"

var (
	ErrorUniqueViolationWhileStoring = errors.New("record with this name already exists")
	ErrorNoRowsWhileGetting          = errors.New("record does not exist")
	ErrorNoRowsUpdated               = errors.New("record does not exist")
	ErrorNoRowsWhileListing          = errors.New("records don't exist")
	ErrorWhileScanning               = errors.New("failed to scan rows")
)
