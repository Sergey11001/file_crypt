package errs

import (
	"errors"
)

// Secure specifies whether to obscure messages of unknown errors.
//
//nolint:gochecknoglobals
var Secure bool

// Unknown code indicates an unknown error. It's the default code.
const Unknown = "Unknown"

// ClassIs reports whether the specified error has the specified class.
//
// ClassIs always returns false if called with [Unclassified].
//
// ClassIs panics if err is nil.
func ClassIs(class Class, err error) bool {
	if err == nil {
		panic("errs: nil error")
	}

	if class == Unclassified {
		return false
	}

	if impl := (interface{ ErrorClass() Class })(nil); errors.As(err, &impl) {
		return class == impl.ErrorClass()
	}

	return class == Internal
}

// CodeIs reports whether the specified error has the specified code.
//
// CodeIs always returns false if called with empty code or [Unknown].
//
// CodeIs panics if err is nil.
func CodeIs(code string, err error) bool {
	if err == nil {
		panic("errs: nil error")
	}

	if code == "" || code == Unknown {
		return false
	}

	if impl := (interface{ ErrorCode() string })(nil); errors.As(err, &impl) {
		return code == impl.ErrorCode()
	}

	return false
}

// Parse parses the specified error.
//
// Parse panics if err is nil.
func Parse(err error) (class Class, code, message string, details map[string]string) {
	if err == nil {
		panic("errs: nil error")
	}

	class = Internal
	if impl := (interface{ ErrorClass() Class })(nil); errors.As(err, &impl) {
		class = impl.ErrorClass()
	}

	code = Unknown
	if impl := (interface{ ErrorCode() string })(nil); errors.As(err, &impl) {
		code = impl.ErrorCode()
		if code == "" {
			code = Unknown
		}
	}

	message = "unknown"
	if !Secure || code != Unknown {
		message = err.Error()
	}

	if impl := (interface{ ErrorDetails() map[string]string })(nil); errors.As(err, &impl) {
		details = impl.ErrorDetails()
	}

	return class, code, message, details
}

type errorWithClass struct {
	error

	class Class
}

func (e *errorWithClass) Unwrap() error {
	return e.error
}

func (e *errorWithClass) ErrorClass() Class {
	return e.class
}

type errorWithCode struct {
	error

	code string
}

func (e *errorWithCode) Unwrap() error {
	return e.error
}

func (e *errorWithCode) ErrorCode() string {
	return e.code
}

type errorWithDetails struct {
	error

	details map[string]string
}

func (e *errorWithDetails) Unwrap() error {
	return e.error
}

func (e *errorWithDetails) ErrorDetails() map[string]string {
	return e.details
}
