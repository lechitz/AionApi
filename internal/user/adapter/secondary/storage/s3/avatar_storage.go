// Package s3 provides the S3-backed avatar storage adapter for the user context.
package s3

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	userOutput "github.com/lechitz/aion-api/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AvatarStorage is an S3-compatible implementation for avatar uploads.
type AvatarStorage struct {
	client        *s3.Client
	bucket        string
	prefix        string
	publicBaseURL string
	log           logger.ContextLogger
}

// NewAvatarStorage creates a new S3-backed avatar storage adapter.
func NewAvatarStorage(cfg config.AvatarStorageConfig, log logger.ContextLogger) (*AvatarStorage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if strings.TrimSpace(cfg.S3Endpoint) != "" {
			o.BaseEndpoint = &cfg.S3Endpoint
		}
		o.UsePathStyle = true
	})

	return &AvatarStorage{
		client:        client,
		bucket:        cfg.S3Bucket,
		prefix:        strings.Trim(cfg.S3Prefix, "/"),
		publicBaseURL: strings.TrimRight(cfg.PublicBaseURL, "/"),
		log:           log,
	}, nil
}

// UploadAvatar uploads avatar object and returns the public URL.
func (s *AvatarStorage) UploadAvatar(ctx context.Context, input userOutput.AvatarUploadInput) (string, error) {
	ctx, span := otel.Tracer("user.adapter.secondary.storage.s3").Start(ctx, "user.avatar_storage.upload")
	defer span.End()

	objectKey := strings.TrimLeft(input.ObjectKey, "/")
	if s.prefix != "" {
		objectKey = s.prefix + "/" + objectKey
	}
	span.SetAttributes(
		attribute.String("avatar.bucket", s.bucket),
		attribute.String("avatar.object_key", objectKey),
		attribute.String("avatar.content_type", input.ContentType),
		attribute.Int("avatar.payload_bytes", len(input.Body)),
	)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &objectKey,
		Body:        bytes.NewReader(input.Body),
		ContentType: &input.ContentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "avatar_upload_failed")
		if s.log != nil {
			s.log.ErrorwCtx(ctx, "failed to upload avatar", "error", err)
		}
		return "", err
	}

	span.SetStatus(codes.Ok, "avatar_uploaded")
	return fmt.Sprintf("%s/%s", s.publicBaseURL, objectKey), nil
}
