package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	mux     *chi.Mux
	handler http.Handler
}

func NewRouter() (*Router, error) {
	r := &Router{}

	r.mux = chi.NewRouter()

	r.mux.Use(
		CORSMiddleware(),
	)

	r.mux.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.mux.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.handler = r.mux

	return r, nil
}

func (r *Router) Mount(part Layouter) {
	if part == nil {
		panic("http server: nil part")
	}

	r.mux.Group(part.Layout)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}
