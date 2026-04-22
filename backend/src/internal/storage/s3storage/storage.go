package s3storage

import (
	"context"
	"fmt"
	"io"

	"github.com/OGZKTeBmj/forum/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client   *s3.Client
	bucket   string
	endpoint string
	base     string
}

func New(ctx context.Context, endpoint, region, bucket, accessKey, secretKey, base string) (*Storage, error) {
	const op = "selectel.New"

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, utils.ErrWrap(op, err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &Storage{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
		base:     base,
	}, nil
}

func (s *Storage) UploadImage(ctx context.Context, key string, body io.Reader, contentType string) error {
	const op = "selectel.UploadImage"

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	return utils.ErrWrap(op, err)
}

func (s *Storage) DeleteImage(ctx context.Context, key string) error {
	const op = "selectel.Delete"

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return utils.ErrWrap(op, err)
}

func (s *Storage) GetPathImage(ctx context.Context) string {
	return s.bucket
}

func (s *Storage) GetPublicUrl(ctx context.Context, key string) string {
	return fmt.Sprintf("%s/%s", s.base, key)
}
