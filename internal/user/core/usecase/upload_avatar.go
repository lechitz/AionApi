package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	userOutput "github.com/lechitz/AionApi/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	maxAvatarUploadSize = int64(20 << 20) // 20MB
)

var allowedAvatarContentTypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/webp": {},
}

// UploadAvatar validates avatar payload, stores it using configured avatar storage, and returns a public avatar URL.
func (s *Service) UploadAvatar(ctx context.Context, cmd input.UploadAvatarCommand) (string, string, int64, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUploadAvatar)
	defer span.End()

	if s.avatarStorage == nil {
		return "", "", 0, fmt.Errorf("avatar storage not configured")
	}

	if cmd.File == nil {
		return "", "", 0, sharederrors.NewValidationError("avatar", "avatar file is required")
	}

	contentType := strings.ToLower(strings.TrimSpace(cmd.ContentType))
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	if contentType != "" {
		if _, ok := allowedAvatarContentTypes[contentType]; !ok {
			return "", "", 0, sharederrors.NewValidationError("avatar", "unsupported image type (allowed: PNG, JPEG, WEBP)")
		}
	}

	if cmd.SizeBytes <= 0 {
		return "", "", 0, sharederrors.NewValidationError("avatar", "empty avatar file")
	}
	if cmd.SizeBytes > maxAvatarUploadSize {
		return "", "", 0, sharederrors.NewValidationError("avatar", fmt.Sprintf("avatar too large (max %d bytes)", maxAvatarUploadSize))
	}

	payload, err := io.ReadAll(io.LimitReader(cmd.File, maxAvatarUploadSize+1))
	if err != nil {
		span.RecordError(err)
		return "", "", 0, err
	}
	if int64(len(payload)) == 0 {
		return "", "", 0, sharederrors.NewValidationError("avatar", "empty avatar file")
	}
	if int64(len(payload)) > maxAvatarUploadSize {
		return "", "", 0, sharederrors.NewValidationError("avatar", fmt.Sprintf("avatar too large (max %d bytes)", maxAvatarUploadSize))
	}

	detectedType := http.DetectContentType(payload)
	if _, ok := allowedAvatarContentTypes[detectedType]; !ok {
		return "", "", 0, sharederrors.NewValidationError("avatar", "invalid image content")
	}

	objectKey := buildAvatarObjectKey(cmd.Filename, detectedType)
	url, uploadErr := s.avatarStorage.UploadAvatar(ctx, userOutput.AvatarUploadInput{
		ObjectKey:   objectKey,
		ContentType: detectedType,
		Body:        payload,
	})
	if uploadErr != nil {
		span.RecordError(uploadErr)
		return "", "", 0, uploadErr
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String("content_type", detectedType),
		attribute.Int64("size_bytes", int64(len(payload))),
		attribute.String("avatar_object_key", objectKey),
	)
	return url, detectedType, int64(len(payload)), nil
}

func buildAvatarObjectKey(filename, contentType string) string {
	ext := contentTypeToExt(contentType)
	if ext == "" {
		ext = filepath.Ext(filename)
	}
	if ext == "" {
		ext = ".img"
	}
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ext
}

func contentTypeToExt(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
