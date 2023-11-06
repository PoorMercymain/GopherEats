// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: subscription.proto

package subscription

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CreateSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateSubscriptionRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateSubscriptionRequestV1MultiError, or nil if none found.
func (m *CreateSubscriptionRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateSubscriptionRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for BundleId

	if len(errors) > 0 {
		return CreateSubscriptionRequestV1MultiError(errors)
	}

	return nil
}

// CreateSubscriptionRequestV1MultiError is an error wrapping multiple
// validation errors returned by CreateSubscriptionRequestV1.ValidateAll() if
// the designated constraints aren't met.
type CreateSubscriptionRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateSubscriptionRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateSubscriptionRequestV1MultiError) AllErrors() []error { return m }

// CreateSubscriptionRequestV1ValidationError is the validation error returned
// by CreateSubscriptionRequestV1.Validate if the designated constraints
// aren't met.
type CreateSubscriptionRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateSubscriptionRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateSubscriptionRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateSubscriptionRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateSubscriptionRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateSubscriptionRequestV1ValidationError) ErrorName() string {
	return "CreateSubscriptionRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e CreateSubscriptionRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateSubscriptionRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateSubscriptionRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateSubscriptionRequestV1ValidationError{}

// Validate checks the field values on ReadSubscriptionRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadSubscriptionRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadSubscriptionRequestV1MultiError, or nil if none found.
func (m *ReadSubscriptionRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadSubscriptionRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	if len(errors) > 0 {
		return ReadSubscriptionRequestV1MultiError(errors)
	}

	return nil
}

// ReadSubscriptionRequestV1MultiError is an error wrapping multiple validation
// errors returned by ReadSubscriptionRequestV1.ValidateAll() if the
// designated constraints aren't met.
type ReadSubscriptionRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadSubscriptionRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadSubscriptionRequestV1MultiError) AllErrors() []error { return m }

// ReadSubscriptionRequestV1ValidationError is the validation error returned by
// ReadSubscriptionRequestV1.Validate if the designated constraints aren't met.
type ReadSubscriptionRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadSubscriptionRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadSubscriptionRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadSubscriptionRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadSubscriptionRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadSubscriptionRequestV1ValidationError) ErrorName() string {
	return "ReadSubscriptionRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadSubscriptionRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadSubscriptionRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadSubscriptionRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadSubscriptionRequestV1ValidationError{}

// Validate checks the field values on ReadSubscriptionResponseV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadSubscriptionResponseV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadSubscriptionResponseV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadSubscriptionResponseV1MultiError, or nil if none found.
func (m *ReadSubscriptionResponseV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadSubscriptionResponseV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for BundleId

	// no validation rules for IsDeleted

	if len(errors) > 0 {
		return ReadSubscriptionResponseV1MultiError(errors)
	}

	return nil
}

// ReadSubscriptionResponseV1MultiError is an error wrapping multiple
// validation errors returned by ReadSubscriptionResponseV1.ValidateAll() if
// the designated constraints aren't met.
type ReadSubscriptionResponseV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadSubscriptionResponseV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadSubscriptionResponseV1MultiError) AllErrors() []error { return m }

// ReadSubscriptionResponseV1ValidationError is the validation error returned
// by ReadSubscriptionResponseV1.Validate if the designated constraints aren't met.
type ReadSubscriptionResponseV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadSubscriptionResponseV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadSubscriptionResponseV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadSubscriptionResponseV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadSubscriptionResponseV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadSubscriptionResponseV1ValidationError) ErrorName() string {
	return "ReadSubscriptionResponseV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadSubscriptionResponseV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadSubscriptionResponseV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadSubscriptionResponseV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadSubscriptionResponseV1ValidationError{}

// Validate checks the field values on ChangeSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ChangeSubscriptionRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ChangeSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ChangeSubscriptionRequestV1MultiError, or nil if none found.
func (m *ChangeSubscriptionRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ChangeSubscriptionRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for BundleId

	// no validation rules for IsDeleted

	if len(errors) > 0 {
		return ChangeSubscriptionRequestV1MultiError(errors)
	}

	return nil
}

// ChangeSubscriptionRequestV1MultiError is an error wrapping multiple
// validation errors returned by ChangeSubscriptionRequestV1.ValidateAll() if
// the designated constraints aren't met.
type ChangeSubscriptionRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ChangeSubscriptionRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ChangeSubscriptionRequestV1MultiError) AllErrors() []error { return m }

// ChangeSubscriptionRequestV1ValidationError is the validation error returned
// by ChangeSubscriptionRequestV1.Validate if the designated constraints
// aren't met.
type ChangeSubscriptionRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ChangeSubscriptionRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ChangeSubscriptionRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ChangeSubscriptionRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ChangeSubscriptionRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ChangeSubscriptionRequestV1ValidationError) ErrorName() string {
	return "ChangeSubscriptionRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ChangeSubscriptionRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sChangeSubscriptionRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ChangeSubscriptionRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ChangeSubscriptionRequestV1ValidationError{}

// Validate checks the field values on CancelSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CancelSubscriptionRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CancelSubscriptionRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CancelSubscriptionRequestV1MultiError, or nil if none found.
func (m *CancelSubscriptionRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *CancelSubscriptionRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	if len(errors) > 0 {
		return CancelSubscriptionRequestV1MultiError(errors)
	}

	return nil
}

// CancelSubscriptionRequestV1MultiError is an error wrapping multiple
// validation errors returned by CancelSubscriptionRequestV1.ValidateAll() if
// the designated constraints aren't met.
type CancelSubscriptionRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CancelSubscriptionRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CancelSubscriptionRequestV1MultiError) AllErrors() []error { return m }

// CancelSubscriptionRequestV1ValidationError is the validation error returned
// by CancelSubscriptionRequestV1.Validate if the designated constraints
// aren't met.
type CancelSubscriptionRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CancelSubscriptionRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CancelSubscriptionRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CancelSubscriptionRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CancelSubscriptionRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CancelSubscriptionRequestV1ValidationError) ErrorName() string {
	return "CancelSubscriptionRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e CancelSubscriptionRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCancelSubscriptionRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CancelSubscriptionRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CancelSubscriptionRequestV1ValidationError{}

// Validate checks the field values on AddBalanceRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AddBalanceRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AddBalanceRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AddBalanceRequestV1MultiError, or nil if none found.
func (m *AddBalanceRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *AddBalanceRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for Amount

	if len(errors) > 0 {
		return AddBalanceRequestV1MultiError(errors)
	}

	return nil
}

// AddBalanceRequestV1MultiError is an error wrapping multiple validation
// errors returned by AddBalanceRequestV1.ValidateAll() if the designated
// constraints aren't met.
type AddBalanceRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AddBalanceRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AddBalanceRequestV1MultiError) AllErrors() []error { return m }

// AddBalanceRequestV1ValidationError is the validation error returned by
// AddBalanceRequestV1.Validate if the designated constraints aren't met.
type AddBalanceRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddBalanceRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddBalanceRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddBalanceRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddBalanceRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddBalanceRequestV1ValidationError) ErrorName() string {
	return "AddBalanceRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e AddBalanceRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddBalanceRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddBalanceRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddBalanceRequestV1ValidationError{}

// Validate checks the field values on ReadUserDataRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadUserDataRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadUserDataRequestV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadUserDataRequestV1MultiError, or nil if none found.
func (m *ReadUserDataRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadUserDataRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	if len(errors) > 0 {
		return ReadUserDataRequestV1MultiError(errors)
	}

	return nil
}

// ReadUserDataRequestV1MultiError is an error wrapping multiple validation
// errors returned by ReadUserDataRequestV1.ValidateAll() if the designated
// constraints aren't met.
type ReadUserDataRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadUserDataRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadUserDataRequestV1MultiError) AllErrors() []error { return m }

// ReadUserDataRequestV1ValidationError is the validation error returned by
// ReadUserDataRequestV1.Validate if the designated constraints aren't met.
type ReadUserDataRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadUserDataRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadUserDataRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadUserDataRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadUserDataRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadUserDataRequestV1ValidationError) ErrorName() string {
	return "ReadUserDataRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadUserDataRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadUserDataRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadUserDataRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadUserDataRequestV1ValidationError{}

// Validate checks the field values on ReadUserDataResponseV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadUserDataResponseV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadUserDataResponseV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadUserDataResponseV1MultiError, or nil if none found.
func (m *ReadUserDataResponseV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadUserDataResponseV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Address

	// no validation rules for BundleId

	// no validation rules for Balance

	if len(errors) > 0 {
		return ReadUserDataResponseV1MultiError(errors)
	}

	return nil
}

// ReadUserDataResponseV1MultiError is an error wrapping multiple validation
// errors returned by ReadUserDataResponseV1.ValidateAll() if the designated
// constraints aren't met.
type ReadUserDataResponseV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadUserDataResponseV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadUserDataResponseV1MultiError) AllErrors() []error { return m }

// ReadUserDataResponseV1ValidationError is the validation error returned by
// ReadUserDataResponseV1.Validate if the designated constraints aren't met.
type ReadUserDataResponseV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadUserDataResponseV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadUserDataResponseV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadUserDataResponseV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadUserDataResponseV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadUserDataResponseV1ValidationError) ErrorName() string {
	return "ReadUserDataResponseV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadUserDataResponseV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadUserDataResponseV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadUserDataResponseV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadUserDataResponseV1ValidationError{}

// Validate checks the field values on ReadBalanceHistoryRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadBalanceHistoryRequestV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadBalanceHistoryRequestV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadBalanceHistoryRequestV1MultiError, or nil if none found.
func (m *ReadBalanceHistoryRequestV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadBalanceHistoryRequestV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for Page

	if len(errors) > 0 {
		return ReadBalanceHistoryRequestV1MultiError(errors)
	}

	return nil
}

// ReadBalanceHistoryRequestV1MultiError is an error wrapping multiple
// validation errors returned by ReadBalanceHistoryRequestV1.ValidateAll() if
// the designated constraints aren't met.
type ReadBalanceHistoryRequestV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadBalanceHistoryRequestV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadBalanceHistoryRequestV1MultiError) AllErrors() []error { return m }

// ReadBalanceHistoryRequestV1ValidationError is the validation error returned
// by ReadBalanceHistoryRequestV1.Validate if the designated constraints
// aren't met.
type ReadBalanceHistoryRequestV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadBalanceHistoryRequestV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadBalanceHistoryRequestV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadBalanceHistoryRequestV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadBalanceHistoryRequestV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadBalanceHistoryRequestV1ValidationError) ErrorName() string {
	return "ReadBalanceHistoryRequestV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadBalanceHistoryRequestV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadBalanceHistoryRequestV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadBalanceHistoryRequestV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadBalanceHistoryRequestV1ValidationError{}

// Validate checks the field values on ReadBalanceHistoryResponseV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ReadBalanceHistoryResponseV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadBalanceHistoryResponseV1 with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ReadBalanceHistoryResponseV1MultiError, or nil if none found.
func (m *ReadBalanceHistoryResponseV1) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadBalanceHistoryResponseV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetHistory() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ReadBalanceHistoryResponseV1ValidationError{
						field:  fmt.Sprintf("History[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ReadBalanceHistoryResponseV1ValidationError{
						field:  fmt.Sprintf("History[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ReadBalanceHistoryResponseV1ValidationError{
					field:  fmt.Sprintf("History[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ReadBalanceHistoryResponseV1MultiError(errors)
	}

	return nil
}

// ReadBalanceHistoryResponseV1MultiError is an error wrapping multiple
// validation errors returned by ReadBalanceHistoryResponseV1.ValidateAll() if
// the designated constraints aren't met.
type ReadBalanceHistoryResponseV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadBalanceHistoryResponseV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadBalanceHistoryResponseV1MultiError) AllErrors() []error { return m }

// ReadBalanceHistoryResponseV1ValidationError is the validation error returned
// by ReadBalanceHistoryResponseV1.Validate if the designated constraints
// aren't met.
type ReadBalanceHistoryResponseV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadBalanceHistoryResponseV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadBalanceHistoryResponseV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadBalanceHistoryResponseV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadBalanceHistoryResponseV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadBalanceHistoryResponseV1ValidationError) ErrorName() string {
	return "ReadBalanceHistoryResponseV1ValidationError"
}

// Error satisfies the builtin error interface
func (e ReadBalanceHistoryResponseV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadBalanceHistoryResponseV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadBalanceHistoryResponseV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadBalanceHistoryResponseV1ValidationError{}

// Validate checks the field values on HistoryElementV1 with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *HistoryElementV1) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HistoryElementV1 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// HistoryElementV1MultiError, or nil if none found.
func (m *HistoryElementV1) ValidateAll() error {
	return m.validate(true)
}

func (m *HistoryElementV1) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for Amount

	// no validation rules for Operation

	// no validation rules for MadeAt

	if len(errors) > 0 {
		return HistoryElementV1MultiError(errors)
	}

	return nil
}

// HistoryElementV1MultiError is an error wrapping multiple validation errors
// returned by HistoryElementV1.ValidateAll() if the designated constraints
// aren't met.
type HistoryElementV1MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HistoryElementV1MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HistoryElementV1MultiError) AllErrors() []error { return m }

// HistoryElementV1ValidationError is the validation error returned by
// HistoryElementV1.Validate if the designated constraints aren't met.
type HistoryElementV1ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HistoryElementV1ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HistoryElementV1ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HistoryElementV1ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HistoryElementV1ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HistoryElementV1ValidationError) ErrorName() string { return "HistoryElementV1ValidationError" }

// Error satisfies the builtin error interface
func (e HistoryElementV1ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHistoryElementV1.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HistoryElementV1ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HistoryElementV1ValidationError{}
