package httpengine

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"univer/internal/dto"
	"univer/internal/httpengine/openapi"
	"univer/pkg/lib/errs"
	"univer/pkg/lib/httpapi"
	"univer/pkg/lib/httpserver"
)

func (c *controller) Refresh(w http.ResponseWriter, r *http.Request) {
	c.base(w, r, httpapi.HandlerWithInput(func(ctx context.Context, input openapi.RefreshInput) (any, error) {
		tokens, err := c.usersService.RefreshTokens(ctx, input.Token)
		if err != nil {
			return nil, err
		}

		return openapi.RefreshResult{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}, nil
	}))
}

func (c *controller) SignIn(w http.ResponseWriter, r *http.Request) {
	c.base(w, r, httpapi.HandlerWithInput(func(ctx context.Context, input openapi.SignInInput) (any, error) {
		tokens, key, err := c.usersService.SignIn(ctx, dto.SignInInput{
			Email:    input.Email,
			Password: input.Password,
		})
		if err != nil {
			return nil, err
		}

		return openapi.SignInResult{
			PublicKey:    key,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			Email:        input.Email,
		}, nil
	}))
}

func (c *controller) SignUp(w http.ResponseWriter, r *http.Request) {
	c.base(w, r, httpapi.HandlerWithInput(func(ctx context.Context, input openapi.SignUpInput) (any, error) {
		tokens, key, err := c.usersService.SignUp(ctx, dto.SignUpInput{
			Name:      input.Name,
			Email:     input.Email,
			Password:  input.Password,
			PublicKey: []byte(input.PublicKey),
		})
		if err != nil {
			return nil, err
		}

		return openapi.SignUpResult{
			PublicKey:    key,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}, nil
	}))
}

func authMiddleware(usersService UsersService, tokenManager TokenManager) httpserver.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			err := func() error {
				tokenString := r.Header.Get("Authorization")
				if tokenString == "" {
					return errs.Unauthenticated.New("NoToken", "no token")
				}

				tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

				userUUID, err := tokenManager.Parse(tokenString)
				if err != nil {
					return err
				}

				id, err := uuid.Parse(userUUID)
				if err != nil {
					return err
				}

				_, err = usersService.User(ctx, id)
				if err != nil {
					return err
				}

				ctx = contextWithUserUUID(ctx, id)

				return nil
			}()
			if err != nil {
				httpapi.WriteError(w, err)

				return
			}

			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type contextUserUUIDKey struct{}

func contextWithUserUUID(parent context.Context, userUUID uuid.UUID) context.Context {
	return context.WithValue(parent, contextUserUUIDKey{}, userUUID)
}

func contextUserUUID(ctx context.Context) (uuid.UUID, bool) {
	var zero uuid.UUID

	val := ctx.Value(contextUserUUIDKey{})
	if val == nil {
		return zero, false
	}

	return val.(uuid.UUID), true
}
