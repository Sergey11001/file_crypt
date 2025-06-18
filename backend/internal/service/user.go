package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"univer/internal/dto"
	"univer/internal/pgrepository"
	"univer/pkg/lib/errs"
)

type UsersConfig struct {
	AccessTokenTTL  time.Duration `json:"access_token_ttl" default:"24h"`
	RefreshTokenTTL time.Duration `json:"refresh_token_ttl" default:"48h"`
}

type UsersService struct {
	config       UsersConfig
	pgClient     PgClient
	redisClient  RedisClient
	hasher       Hasher
	tokenManager TokenManager
}

func NewUsersService(
	config UsersConfig,
	pgClient PgClient,
	redisClient RedisClient,
	hasher Hasher,
	tokenManager TokenManager,
) *UsersService {
	return &UsersService{
		config:       config,
		pgClient:     pgClient,
		redisClient:  redisClient,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input dto.SignUpInput) (dto.Tokens, string, error) {
	var zero dto.Tokens

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return zero, "", err
	}

	user, err := pgrepository.New(s.pgClient).CreateUser(ctx, pgrepository.CreateUserParams{
		Uuid:         uuid.New(),
		PasswordHash: string(passwordHash),
		Email:        input.Email,
		Name:         input.Name,
		PublicKey:    input.PublicKey,
	})
	if err != nil {
		if pgrepository.IsUserEmailKeyViolation(err) {
			return zero, "", errs.Invalid.New("EmailAlreadyExists", "email already exists")
		}

		return zero, "", err
	}

	tokens, err := s.createSession(ctx, user.Uuid)
	if err != nil {
		return zero, "", err
	}

	return tokens, string(user.PublicKey), nil
}

func (s *UsersService) SignIn(ctx context.Context, input dto.SignInInput) (dto.Tokens, string, error) {
	var zero dto.Tokens

	user, err := pgrepository.New(s.pgClient).UserByEmail(ctx, pgrepository.UserByEmailParams{
		Email: input.Email,
	})
	if err != nil {
		if pgrepository.IsNoRows(err) {
			return zero, "", errs.NotFound.New("UserNotFound", "user not found")
		}

		return zero, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return zero, "", errs.PermissionDenied.New("PermissionDenied", "permission denied")
	}

	tokens, err := s.createSession(ctx, user.Uuid)
	if err != nil {
		return zero, "", err
	}

	return tokens, string(user.PublicKey), nil
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (dto.Tokens, error) {
	var zero dto.Tokens

	val, err := s.redisClient.Get(ctx, refreshToken).Result()
	if err != nil {
		return zero, fmt.Errorf("failed to get value from Redis: %v", err)
	}

	uuid, err := uuid.Parse(val)
	if err != nil {
		return zero, fmt.Errorf("failed to parse UUID: %v", err)
	}

	return s.createSession(ctx, uuid)
}

func (s *UsersService) createSession(ctx context.Context, userUUID uuid.UUID) (dto.Tokens, error) {
	var (
		res dto.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userUUID.String(), s.config.AccessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	err = s.redisClient.Set(ctx, res.RefreshToken, userUUID, s.config.RefreshTokenTTL).Err()
	if err != nil {
		return res, err
	}

	return res, err
}

func (s *UsersService) User(ctx context.Context, userUUID uuid.UUID) (dto.User, error) {
	var zero dto.User

	user, err := pgrepository.New(s.pgClient).User(ctx, pgrepository.UserParams{
		Uuid: userUUID,
	})
	if err != nil {
		if pgrepository.IsNoRows(err) {
			return zero, errs.NotFound.New("UserNotFound", "user not found")
		}

		return zero, err
	}

	return dto.User{
		UUID:      userUUID,
		Email:     user.Email,
		Name:      user.Name,
		PublicKey: string(user.PublicKey),
	}, nil
}

func (s *UsersService) Users(ctx context.Context, userUUID uuid.UUID) ([]dto.User, error) {
	users, err := pgrepository.New(s.pgClient).Users(ctx, pgrepository.UsersParams{
		Uuid: userUUID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, dto.User{
			UUID:      user.Uuid,
			Email:     user.Email,
			Name:      user.Name,
			PublicKey: string(user.PublicKey),
		})
	}

	return result, err
}

func (s *UsersService) UsersForShare(ctx context.Context, fileUUID uuid.UUID) ([]dto.User, error) {
	users, err := pgrepository.New(s.pgClient).UsersForShare(ctx, pgrepository.UsersForShareParams{
		FileUuid: fileUUID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, dto.User{
			UUID:      user.Uuid,
			Email:     user.Email,
			Name:      user.Name,
			PublicKey: string(user.PublicKey),
		})
	}

	return result, err
}

func (s *UsersService) AvailableUsers(ctx context.Context, userUUID, fileUUID uuid.UUID) ([]dto.User, error) {
	users, err := pgrepository.New(s.pgClient).AvailableUsers(ctx, pgrepository.AvailableUsersParams{
		Uuid:     userUUID,
		FileUuid: fileUUID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]dto.User, 0, len(users))
	for _, user := range users {
		result = append(result, dto.User{
			UUID:      user.Uuid,
			Email:     user.Email,
			Name:      user.Name,
			PublicKey: string(user.PublicKey),
		})
	}

	return result, err
}
