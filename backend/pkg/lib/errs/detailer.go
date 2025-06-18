package errs

import (
	"errors"
)

// Detailer modifies error details.
type Detailer interface {
	DetailError(details map[string]string)
}

// DetailerFunc is the functional implementation of the [Detailer] interface.
type DetailerFunc func(details map[string]string)

// DetailError implements the [Detailer] interface.
func (f DetailerFunc) DetailError(details map[string]string) {
	f(details)
}

// Details is [Detailer] setting details with keys and values from map.
type Details map[string]string

// DetailError implements the [Detailer] interface.
func (d Details) DetailError(details map[string]string) {
	for k, v := range d {
		details[k] = v
	}
}

// Detail returns a new error wrapping the specified one with details.
//
// Detail returns err as is if no detailers are specified. Otherwise, it considers existing details, if any.
//
// Detail panics if err or detailer are nil.
func Detail(err error, detailers ...Detailer) error {
	if err == nil {
		panic("errs: nil error")
	}

	if len(detailers) == 0 {
		return err
	}

	details := make(map[string]string)

	if impl := (interface{ ErrorDetails() map[string]string })(nil); errors.As(err, &impl) {
		for k, v := range impl.ErrorDetails() {
			details[k] = v
		}
	}

	for _, detailer := range detailers {
		if detailer == nil {
			panic("errs: nil detailer")
		}

		detailer.DetailError(details)
	}

	if len(details) == 0 {
		details = nil
	}

	return &errorWithDetails{err, details}
}
