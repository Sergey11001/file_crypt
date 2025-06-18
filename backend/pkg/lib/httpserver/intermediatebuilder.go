package httpserver

import (
	"context"
	"net/http"
)

type IntermediateBuilder struct {
	middlewares []Middleware
}

func NewIntermediateBuilder() *IntermediateBuilder {
	return &IntermediateBuilder{}
}

func (b *IntermediateBuilder) Use(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		if middleware == nil {
			panic("http server: nil middleware")
		}
	}

	b.middlewares = append(b.middlewares, middlewares...)
}

func (b *IntermediateBuilder) Copy() *IntermediateBuilder {
	bc := &IntermediateBuilder{
		middlewares: make([]Middleware, len(b.middlewares)),
	}

	copy(bc.middlewares, b.middlewares)

	return bc
}

func (b *IntermediateBuilder) Build() IntermediateFunc {
	contextHandlerKey := new(int)

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w == nil {
			panic("http server: nil response writer")
		}
		if r == nil {
			panic("http server: nil request")
		}

		val := r.Context().Value(contextHandlerKey)
		if val == nil {
			panic("http server: no context handler")
		}

		val.(http.Handler).ServeHTTP(w, r)
	})

	for i := len(b.middlewares) - 1; i >= 0; i-- {
		h = b.middlewares[i].Wrap(h)
	}

	return func(w http.ResponseWriter, r *http.Request, handler http.Handler) {
		if w == nil {
			panic("http server: nil response writer")
		}
		if r == nil {
			panic("http server: nil request")
		}
		if handler == nil {
			panic("http server: nil handler")
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextHandlerKey, handler)))
	}
}
