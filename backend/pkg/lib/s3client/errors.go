package s3client

import (
	"errors"
)

var ErrDownloadingFailed = errors.New("s3 client: downloading failed")

var ErrGettingBucketsFailed = errors.New("s3 client: getting buckets failed")

var ErrInitializationFailed = errors.New("s3 client: initialization failed")

var ErrMakingBucketFailed = errors.New("s3 client: making bucket failed")

var ErrNoAccessKeyID = errors.New("s3 client: no access key id")

var ErrNoEndpoint = errors.New("s3 client: no endpoint")

var ErrNoSecretAccessKey = errors.New("s3 client: no secret access key")

var ErrReadingFailed = errors.New("s3 client: reading failed")

var ErrUploadingFailed = errors.New("s3 client: uploading failed")

var ErrDeletingFailed = errors.New("s3 client: deleting failed")
