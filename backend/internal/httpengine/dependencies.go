package httpengine

import (
	"context"

	"github.com/google/uuid"

	"univer/internal/dto"
)

type FilesService interface {
	CreateFile(ctx context.Context, userUUID uuid.UUID, input dto.CreateFileInput) (dto.File, error)
	DeleteFile(ctx context.Context, userUUID, fileUUID uuid.UUID) error
	UpdateUserKeys(ctx context.Context, userUUID uuid.UUID) error
	File(ctx context.Context, userUUID, fileUUID uuid.UUID) (string, []byte, error)
	CommonFile(ctx context.Context, fileUUID uuid.UUID) (string, []byte, error)
	Files(ctx context.Context, userUUID uuid.UUID) ([]dto.File, error)
	ShareFile(ctx context.Context, userUUID, fileUUID uuid.UUID, input dto.ShareFileInput) error
	DeleteFileAccess(ctx context.Context, input dto.DeleteFileAccessInput) error
	AvailableFiles(ctx context.Context, userUUID uuid.UUID) ([]dto.File, error)
}

type UsersService interface {
	SignUp(ctx context.Context, input dto.SignUpInput) (dto.Tokens, string, error)
	SignIn(ctx context.Context, input dto.SignInInput) (dto.Tokens, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (dto.Tokens, error)
	AvailableUsers(ctx context.Context, userUUID, fileUUID uuid.UUID) ([]dto.User, error)
	User(ctx context.Context, userUUID uuid.UUID) (dto.User, error)
	Users(ctx context.Context, userUUID uuid.UUID) ([]dto.User, error)
	UsersForShare(ctx context.Context, fileUUID uuid.UUID) ([]dto.User, error)
}

type TokenManager interface {
	Parse(accessToken string) (string, error)
}

type S3Client interface {
	Download(ctx context.Context, bucket, path string) ([]byte, error)
	Upload(ctx context.Context, bucket, path string, data []byte) error
}

type Logger interface {
	Error(format string, args ...interface{})
}
