package file_gateway

import (
	"context"
	"io"

	"finance-backend/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Gateway struct {
	s3Client *s3.S3
	log      *logger.Logger
}

func NewS3Gateway(sess *session.Session, log *logger.Logger) *S3Gateway {
	return &S3Gateway{
		s3Client: s3.New(sess),
		log:      log,
	}
}

func (g *S3Gateway) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := g.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		g.log.Error(ctx, "s3_operation_error", map[string]interface{}{"error": err.Error(), "bucket": bucket, "key": key})
		return nil, 0, err
	}
	return result.Body, aws.Int64Value(result.ContentLength), nil
}

func (g *S3Gateway) UploadObject(ctx context.Context, bucket, key string, body io.ReadSeeker, size int64, contentType string) (*UploadResult, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	}
	result, err := g.s3Client.PutObjectWithContext(ctx, input)
	if err != nil {
		g.log.Error(ctx, "s3_operation_error", map[string]interface{}{"error": err.Error(), "bucket": bucket, "key": key})
		return nil, err
	}
	return &UploadResult{
		Bucket: bucket,
		Key:    key,
		Size:   size,
		ETag:   aws.StringValue(result.ETag),
	}, nil
}

func (g *S3Gateway) DeleteObject(ctx context.Context, bucket, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err := g.s3Client.DeleteObjectWithContext(ctx, input)
	return err
}
