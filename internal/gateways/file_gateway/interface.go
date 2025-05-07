package file_gateway

import (
	"context"
	"io"
)

type IFileGateway interface {
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, int64, error)

	UploadObject(ctx context.Context, bucket, key string, body io.ReadSeeker, size int64, contentType string) (*UploadResult, error)

	DeleteObject(ctx context.Context, bucket, key string) error
}

type UploadResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Size   int64  `json:"size"`
	ETag   string `json:"etag"`
}
