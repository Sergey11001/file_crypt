package errs

import (
	"net/http"
)

// HTTPMethodNotAllowed class is intended only for compatibility with HTTP status 405 Method Not Allowed.
//
// Don't use this class unless you know exactly what you're doing.
const HTTPMethodNotAllowed Class = "HTTPMethodNotAllowed"

//nolint:gochecknoglobals
var httpClasses = map[int]Class{
	http.StatusBadRequest:          Invalid,
	http.StatusUnauthorized:        Unauthenticated,
	http.StatusForbidden:           PermissionDenied,
	http.StatusNotFound:            NotFound,
	http.StatusMethodNotAllowed:    HTTPMethodNotAllowed,
	http.StatusInternalServerError: Internal,
	http.StatusNotImplemented:      Unimplemented,
	http.StatusServiceUnavailable:  Unavailable,
}

// HTTPClass converts the specified HTTP status code to class.
//
// HTTPClass converts any unknown code to [Unclassified].
func HTTPClass(statusCode int) Class {
	class, ok := httpClasses[statusCode]
	if !ok {
		class = Unclassified
	}

	return class
}

//nolint:gochecknoglobals
var httpStatuses = map[Class]int{
	Invalid:              http.StatusBadRequest,
	Unauthenticated:      http.StatusUnauthorized,
	PermissionDenied:     http.StatusForbidden,
	NotFound:             http.StatusNotFound,
	HTTPMethodNotAllowed: http.StatusMethodNotAllowed,
	Internal:             http.StatusInternalServerError,
	Unimplemented:        http.StatusNotImplemented,
	Unavailable:          http.StatusServiceUnavailable,
}

// HTTPStatus converts the specified class to HTTP status code.
//
// HTTPStatus converts any unknown class to 500 Internal Server Error.
func HTTPStatus(class Class) int {
	statusCode, ok := httpStatuses[class]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	return statusCode
}
