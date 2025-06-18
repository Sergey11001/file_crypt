package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Intermediate interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, handler http.Handler)
}

type IntermediateFunc func(w http.ResponseWriter, r *http.Request, handler http.Handler)

func (f IntermediateFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, handler http.Handler) {
	f(w, r, handler)
}

type Layouter interface {
	Layout(mux chi.Router)
}

type LayouterFunc func(mux chi.Router)

func (f LayouterFunc) Layout(mux chi.Router) {
	f(mux)
}

type Middleware interface {
	Wrap(handler http.Handler) http.Handler
}

type MiddlewareFunc func(handler http.Handler) http.Handler

func (f MiddlewareFunc) Wrap(handler http.Handler) http.Handler {
	return f(handler)
}
