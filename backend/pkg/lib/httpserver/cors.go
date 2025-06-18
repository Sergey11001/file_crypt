package httpserver

import (
	"net/http"

	"github.com/go-chi/cors"
)

func CORSMiddleware() MiddlewareFunc {
	return cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
		AllowedOrigins:     []string{"*"},
	})
}
