package httpengine

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"univer/internal/httpengine/openapi"
	"univer/pkg/lib/errs"
	"univer/pkg/lib/httpapi"
	"univer/pkg/lib/httpserver"
)

func New(
	logger Logger,
	filesService FilesService,
	usersService UsersService,
	tokenManager TokenManager,
) (httpserver.LayouterFunc, error) {
	if logger == nil {
		panic("logger is nil")
	}
	if filesService == nil {
		panic("files service is nil")
	}
	if usersService == nil {
		panic("users service is nil")
	}
	if tokenManager == nil {
		panic("token manager is nil")
	}

	return func(r chi.Router) {
		base := httpserver.NewIntermediateBuilder()

		auth := base.Copy()
		auth.Use(authMiddleware(usersService, tokenManager))

		r.Use(
			middleware.Recoverer,
		)

		controller := &controller{
			logger,
			filesService,
			usersService,

			base.Build(),
			auth.Build(),
		}

		r.Route("/api", func(apiRouter chi.Router) {
			openapi.HandlerWithOptions(controller, openapi.ChiServerOptions{
				BaseRouter: apiRouter,
				ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					httpapi.WriteError(w, errs.Invalid.As("BadRequest", err))
				},
			})
		})
	}, nil
}

type controller struct {
	logger       Logger
	filesService FilesService
	usersService UsersService

	base httpserver.IntermediateFunc
	auth httpserver.IntermediateFunc
}
