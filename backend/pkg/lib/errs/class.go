package errs

import (
	"errors"
)

// Class of errors.
type Class string

// Unclassified class indicated error that can't be classified.
//
// Unclassified class is converted to HTTP status 500 Internal Server Error.
//
// Unclassified class is converted to gRPC status 2 UNKNOWN.
const Unclassified Class = "Unclassified"

// Internal class indicates internal error. It's the default class.
//
// Internal class corresponds to HTTP status 500 Internal Server Error.
//
// Internal class corresponds to gRPC status 13 INTERNAL.
const Internal Class = "Internal"

// Invalid class indicates that request is invalid.
//
// Invalid class corresponds to HTTP status 400 Bad Request.
//
// Invalid class corresponds to gRPC status 3 INVALID_ARGUMENT.
const Invalid Class = "Invalid"

// NotFound class indicates that requested entity is not found.
//
// NotFound class corresponds to HTTP status 404 Not Found.
//
// NotFound class corresponds to gRPC status 5 NOT_FOUND.
const NotFound Class = "NotFound"

// PermissionDenied class indicates that caller doesn't have permission to execute operation.
//
// PermissionDenied class corresponds to HTTP status 403 Forbidden.
//
// PermissionDenied class corresponds to gRPC status 7 PERMISSION_DENIED.
const PermissionDenied Class = "PermissionDenied"

// Unauthenticated class indicates that request doesn't have valid authentication credentials for operation.
//
// Unauthenticated class corresponds to HTTP status 401 Unauthorized.
//
// Unauthenticated class corresponds to gRPC status 16 UNAUTHENTICATED.
const Unauthenticated Class = "Unauthenticated"

// Unavailable class indicates that service is currently unavailable.
//
// Unavailable class corresponds to HTTP status 503 Service Unavailable.
//
// Unavailable class corresponds to gRPC status 14 UNAVAILABLE.
const Unavailable Class = "Unavailable"

// Unimplemented class indicates that operation isn't implemented or isn't supported/enabled.
//
// Unimplemented class corresponds to HTTP status 501 Not Implemented.
//
// Unimplemented class corresponds to gRPC status 12 UNIMPLEMENTED.
const Unimplemented Class = "Unimplemented"

// New returns a new error of this class with the specified code, message and details.
//
// New panics if detailer is nil.
func (c Class) New(code, message string, detailers ...Detailer) error {
	return c.As(code, errors.New(message), detailers...) //nolint:goerr113
}

// As returns a new error of this class wrapping the specified one with the specified code and details.
//
// As panics if err or detailer are nil.
func (c Class) As(code string, err error, detailers ...Detailer) error {
	if err == nil {
		panic("errs: nil error")
	}

	return c.Cast(&errorWithCode{err, code}, detailers...)
}

// Cast returns a new error of this class wrapping the specified one with the specified details.
//
// Cast panics if err or detailer are nil.
func (c Class) Cast(err error, detailers ...Detailer) error {
	if err == nil {
		panic("errs: nil error")
	}

	return Detail(&errorWithClass{err, c}, detailers...)
}
