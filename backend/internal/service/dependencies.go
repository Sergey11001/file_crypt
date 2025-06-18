package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"univer/internal/pgrepository"
)

type S3Client interface {
	Delete(ctx context.Context, bucket, path string) error
	Download(ctx context.Context, bucket, path string) ([]byte, error)
	Upload(ctx context.Context, bucket, path string, data []byte) error
}

type PgClient interface {
	pgrepository.DBTX
}

type RedisClient interface {
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type Hasher interface {
	Hash(password string) (string, error)
}

type TokenManager interface {
	NewJWT(userUUID string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}
