package s3client

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"univer/pkg/lib/types"
)

type S3Client struct {
	logger Logger
	client *minio.Client
}

type ClientConfig struct {
	AccessKeyID     string `default:""`
	SecretAccessKey string `default:""`
	Region          string `default:""`
	Endpoint        string `default:""`
	Secure          bool   `default:"false"`
}

func New(config ClientConfig, logger Logger) (*S3Client, error) {
	if logger == nil {
		panic("s3 client: nil logger")
	}

	if config.AccessKeyID == "" {
		return nil, ErrNoAccessKeyID
	}
	if config.SecretAccessKey == "" {
		return nil, ErrNoSecretAccessKey
	}
	if config.Endpoint == "" {
		return nil, ErrNoEndpoint
	}

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.Secure,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInitializationFailed, err)
	}

	return &S3Client{
		logger: logger,
		client: client,
	}, nil
}

func (c *S3Client) Upload(ctx context.Context, bucket, path string, data []byte) error {
	if ctx == nil {
		panic("s3 client: nil context")
	}

	_, err := c.client.PutObject(
		ctx,
		bucket,
		path,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUploadingFailed, err)
	}

	return nil
}

func (c *S3Client) Delete(ctx context.Context, bucket, path string) error {
	if ctx == nil {
		panic("s3 client: nil context")
	}

	err := c.client.RemoveObject(ctx, bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeletingFailed, err)
	}

	return nil
}

func (c *S3Client) Download(ctx context.Context, bucket, path string) ([]byte, error) {
	if ctx == nil {
		panic("s3 client: nil context")
	}

	object, err := c.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDownloadingFailed, err)
	}
	defer func() {
		_ = object.Close()
	}()

	var buf bytes.Buffer

	if _, err = io.Copy(&buf, object); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadingFailed, err)
	}

	return buf.Bytes(), nil
}

func (c *S3Client) Stat(ctx context.Context, bucket, key string) bool {
	if ctx == nil {
		panic("s3 client: nil context")
	}

	_, err := c.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})

	return err == nil
}

func (c *S3Client) InitBuckets(buckets []string) error {
	ctx := context.Background()

	bucketList, err := c.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGettingBucketsFailed, err)
	}

	existingBuckets := make([]string, 0, len(bucketList))
	for _, bucket := range bucketList {
		existingBuckets = append(existingBuckets, bucket.Name)
	}

	notExistingBuckets := types.SliceDiff(buckets, existingBuckets)
	for _, bucket := range notExistingBuckets {
		err = c.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%w: %w", ErrMakingBucketFailed, err)
		}
	}

	return nil
}
