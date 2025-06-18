package httpengine

import (
	"context"
	"net/http"

	"github.com/oapi-codegen/runtime/types"

	"univer/internal/httpengine/openapi"
	"univer/pkg/lib/errs"
	"univer/pkg/lib/httpapi"
)

func (c *controller) AvailableUsers(w http.ResponseWriter, r *http.Request, fileUUID types.UUID) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		users, err := c.usersService.AvailableUsers(ctx, userUUID, fileUUID)
		if err != nil {
			return nil, err
		}

		return openapi.UsersResult{
			Users: users,
		}, nil
	}))
}

func (c *controller) Users(w http.ResponseWriter, r *http.Request) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		userUUID, ok := contextUserUUID(ctx)
		if !ok {
			return nil, errs.PermissionDenied.New("Unauthorized", "unauthorized")
		}

		users, err := c.usersService.Users(ctx, userUUID)
		if err != nil {
			return nil, err
		}

		return openapi.UsersResult{
			Users: users,
		}, nil
	}))
}

func (c *controller) UsersForShare(w http.ResponseWriter, r *http.Request, fileUUID types.UUID) {
	c.auth(w, r, httpapi.Handler(func(ctx context.Context) (any, error) {
		users, err := c.usersService.UsersForShare(ctx, fileUUID)
		if err != nil {
			return nil, err
		}

		return openapi.UsersResult{
			Users: users,
		}, nil
	}))
}
