package errs

import (
	"errors"
	"fmt"
)

// Unhandled class to catch unhandled errors.
const Unhandled Class = "Unhandled"

type DetailsProvider interface {
	Details() map[string]string
}

type Marker[M any] interface {
	comparable
	Mark() M
}

// G is a generic error. it holds mark interface and knows its origin pointer.
type G[M Marker[M]] struct {
	reason  string
	origin  *G[M]
	class   Class
	details map[string]string
}

func (b *G[M]) Mark() M {
	var m M

	return m
}

func (b *G[M]) Details() map[string]string {
	return b.details
}

func mergeMaps(a, b map[string]string) map[string]string {
	c := make(map[string]string, len(a)+len(b))

	for k, v := range a {
		c[k] = v
	}

	for k, v := range b {
		c[k] = v
	}

	return c
}

func (b *G[M]) WithDetails(m map[string]string) *G[M] {
	return &G[M]{
		origin:  b,
		class:   b.class,
		details: mergeMaps(b.details, m),
	}
}

func (b *G[M]) ErrorClass() Class {
	class := b.Origin().class

	if class == "" {
		return Unclassified
	}

	return class
}

func (b *G[M]) Error() string {
	if b.origin != nil {
		str := b.origin.Error()
		if b.details != nil || b.reason == "" {
			return str
		}
		if str != "" {
			return fmt.Sprintf("%s: %s", str, b.reason)
		}
	}

	return b.reason
}

func (b *G[M]) Origin() *G[M] {
	if b.origin != nil {
		return b.origin.Origin()
	}

	return b
}

// Wrap adds more clarification to error.
func (b *G[M]) Wrap(e error) *G[M] {
	reason := "NIL WRAP ERROR"
	if e != nil {
		reason = e.Error()
	}

	details := b.details
	if g := (*G[M])(nil); errors.As(e, &g) {
		details = mergeMaps(b.details, g.details)
	}

	return &G[M]{
		reason:  reason,
		origin:  b,
		details: details,
	}
}

// WithReason adds reason to error.
func (b *G[M]) WithReason(reason string) *G[M] {
	return &G[M]{
		reason:  reason,
		origin:  b,
		details: b.details,
	}
}

// EnsureG makes sure error is generic marked M.
func EnsureG[M Marker[M]](e error) *G[M] {
	if e == nil {
		return &G[M]{
			reason: "NIL ENSURE ERROR",
		}
	}

	if g := (*G[M])(nil); errors.As(e, &g) {
		return g
	}

	return NewUnhandled[M](fmt.Sprintf("unhandled in ensure %T: %s", e, e.Error()))
}

func NewInternal[M Marker[M]]() *G[M] {
	return &G[M]{
		class: Internal,
	}
}

func NewInvalid[M Marker[M]]() *G[M] {
	return &G[M]{
		class: Invalid,
	}
}

func NewNotFound[M Marker[M]](reason string) *G[M] {
	return &G[M]{
		reason: reason,
		class:  NotFound,
	}
}

func NewUnclassified[M Marker[M]](reason string) *G[M] {
	return &G[M]{
		reason: reason,
		class:  Unclassified,
	}
}

func NewUnhandled[M Marker[M]](reason string) *G[M] {
	return &G[M]{
		reason: reason,
		class:  Unhandled,
	}
}

type Map[A Marker[A], B Marker[B]] map[*G[A]]*G[B]

func Mapper[A Marker[A], B Marker[B]](err error, m map[*G[A]]*G[B]) *G[B] {
	g := EnsureG[A](err)

	e, ok := m[g.Origin()]
	if ok {
		return e.WithDetails(g.details)
	}

	return &G[B]{
		reason:  fmt.Sprintf("unhandled in mapper %T: %s", err, g.reason),
		class:   g.class,
		details: g.details,
	}
}
