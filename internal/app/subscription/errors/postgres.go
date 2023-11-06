package errors

import "errors"

var (
	ErrorUniqueViolationWhileCreating = errors.New("the user already have a subscription, to change it, use another endpoint")
	ErrorNoRowsWhileReading           = errors.New("the user does not have active subscription or does not exist")
	ErrorNoRowsUpdated                = errors.New("the user does not have active subscription or does not exist")
	ErrorNoRowsWhileReadingHistory    = errors.New("the user does not have operations on the provided page or does not exist")
	ErrorNotEnoughFunds               = errors.New("the user does not have enough funds on balance or does not exist")
)
