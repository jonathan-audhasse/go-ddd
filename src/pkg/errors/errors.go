package errors

import (
	"errors"
	"fmt"
)

const NoID = ErrorID(0)

// Error ID of the error
type ErrorID uint

type customError struct {
	errID   ErrorID
	message string
	cause   error
}

// New creates a new customError
func (et ErrorID) New(msg string) error {
	return customError{errID: et, message: msg}
}

// Newf creates a new customError with formatted message
func (et ErrorID) Newf(msg string, args ...interface{}) error {
	return customError{errID: et, message: fmt.Sprintf(msg, args...)}
}

// Wrap creates a new wrapped error
func (et ErrorID) Wrap(err error, msg string) error {
	return customError{errID: et, message: msg, cause: err}
}

// Wrapf creates a new wrapped error with formatted message
func (et ErrorID) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{errID: et, message: fmt.Sprintf(msg, args...), cause: err}
}

// AsError create a new customer Error
func (et ErrorID) AsError() error {
	return customError{errID: et}
}

// Error returns the message of a customError
func (e customError) Error() string {
	return e.message
}

func (e customError) Unwrap() error {
	return e.cause
}

// Cause returns the cause of a customError
func (e customError) Cause() error {
	return e.cause
}

// ID returns the errID of a customError
func (e customError) ID() ErrorID {
	return e.errID
}

func (e customError) Is(target error) bool {
	if t, ok := target.(customError); ok {
		return e.errID == t.errID
	}
	return false
}

// FullError returns the full unwrapped error of a customError
func (e customError) FullError() string {
	return e.computeFullError(NoID, "")
}

func (e customError) computeFullError(previousErrorID ErrorID, fullErr string) string {
	// display type only if valid and not already displayed
	if e.errID != NoID && e.errID != previousErrorID {
		fullErr = fmt.Sprintf("%s(%v) ", fullErr, e.errID)
	}
	// display message
	fullErr = fmt.Sprintf("%s%s", fullErr, e.message)
	if e.cause == nil {
		return fullErr
	}
	cause, ok := e.cause.(customError)
	if !ok {
		// display underlying cause
		fullErr = fmt.Sprintf("%s: %s", fullErr, e.cause.Error())
		return fullErr
	}
	return cause.computeFullError(e.errID, fullErr+": ")
}

// New creates a no type error
func New(msg string) error {
	return NoID.New(msg)
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return NoID.Newf(msg, args...)
}

// Wrap wrans an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, "%s", msg)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	if customErr, ok := err.(customError); ok {
		return customErr.errID.Wrapf(err, msg, args...)
	}
	return NoID.Wrapf(err, msg, args...)
}

// Unwrap unwrap an error
func Unwrap(err error) error {
	if customErr, ok := err.(customError); ok {
		return customErr.Unwrap()
	}
	return errors.Unwrap(err)
}

// FullError get the full unwrapped error
func FullError(err error) string {
	if customErr, ok := err.(customError); ok {
		return customErr.FullError()
	}
	if err == nil {
		return ""
	}
	return err.Error()
}

// GetID returns the error ID
func GetID(err error) ErrorID {
	if customErr, ok := err.(customError); ok {
		return customErr.errID
	}
	return NoID
}
