package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appconfig "github.com/kiritoxkiriko/comical-tool/server/internal/config"
)

type S3 struct {
	client *s3.Client
	bucket string
}

func NewS3(ctx context.Context, cfg appconfig.StorageConfig) (*S3, error) {
	if cfg.S3Bucket == "" {
		return nil, errors.New("storage.s3_bucket is required")
	}
	region := cfg.S3Region
	if region == "" {
		region = "auto"
	}
	options := []func(*config.LoadOptions) error{config.WithRegion(region)}
	if cfg.S3AccessKeyID != "" || cfg.S3SecretAccessKey != "" {
		options = append(options, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.S3AccessKeyID, cfg.S3SecretAccessKey, ""),
		))
	}
	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(awsCfg, func(opts *s3.Options) {
		opts.UsePathStyle = cfg.S3UsePathStyle
		if cfg.S3Endpoint != "" {
			opts.BaseEndpoint = aws.String(cfg.S3Endpoint)
		}
	})
	return &S3{client: client, bucket: cfg.S3Bucket}, nil
}

func (s *S3) Put(ctx context.Context, key string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (s *S3) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3) Head(ctx context.Context, key string) (ObjectInfo, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return ObjectInfo{}, fmt.Errorf("object %q not found: %w", key, err)
		}
		return ObjectInfo{}, err
	}
	return ObjectInfo{Size: aws.ToInt64(result.ContentLength), ContentType: aws.ToString(result.ContentType)}, nil
}
