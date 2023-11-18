// Package email contains email address handling utilities.
package email

import "net/mail"

// ValidateEmail validates email address and returns false in case it's invalid.
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}
