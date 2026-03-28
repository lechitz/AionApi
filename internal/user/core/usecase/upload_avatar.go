package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	userOutput "github.com/lechitz/aion-api/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	maxAvatarUploadSize = int64(20 << 20) // 20MB
)

// UploadAvatar validates avatar payload, stores it using configured avatar storage, and returns a public avatar URL.
func (s *Service) UploadAvatar(ctx context.Context, cmd input.UploadAvatarCommand) (string, string, int64, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUploadAvatar)
	defer span.End()

	if s.avatarStorage == nil {
		return "", "", 0, errors.New(errAvatarStorageNotConfigured)
	}

	if cmd.File == nil {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, errAvatarFileRequired)
	}

	contentType := strings.ToLower(strings.TrimSpace(cmd.ContentType))
	if idx := strings.Index(contentType, contentTypeParamSeparator); idx >= 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	if contentType != "" {
		if !isAllowedAvatarContentType(contentType) {
			return "", "", 0, sharederrors.NewValidationError(fieldAvatar, errAvatarUnsupportedType)
		}
	}

	if cmd.SizeBytes <= 0 {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, errAvatarEmptyFile)
	}
	if cmd.SizeBytes > maxAvatarUploadSize {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, fmt.Sprintf(errAvatarTooLarge, maxAvatarUploadSize))
	}

	payload, err := io.ReadAll(io.LimitReader(cmd.File, maxAvatarUploadSize+1))
	if err != nil {
		span.RecordError(err)
		return "", "", 0, err
	}
	if int64(len(payload)) == 0 {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, errAvatarEmptyFile)
	}
	if int64(len(payload)) > maxAvatarUploadSize {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, fmt.Sprintf(errAvatarTooLarge, maxAvatarUploadSize))
	}

	detectedType := http.DetectContentType(payload)
	if !isAllowedAvatarContentType(detectedType) {
		return "", "", 0, sharederrors.NewValidationError(fieldAvatar, errAvatarInvalidContent)
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

func isAllowedAvatarContentType(contentType string) bool {
	switch contentType {
	case contentTypeImagePNG, contentTypeImageJPEG, contentTypeImageWEBP:
		return true
	default:
		return false
	}
}

func buildAvatarObjectKey(filename, contentType string) string {
	ext := contentTypeToExt(contentType)
	if ext == "" {
		ext = filepath.Ext(filename)
	}
	if ext == "" {
		ext = defaultAvatarFallbackExt
	}
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ext
}

func contentTypeToExt(contentType string) string {
	switch contentType {
	case contentTypeImagePNG:
		return extPNG
	case contentTypeImageJPEG:
		return extJPEG
	case contentTypeImageWEBP:
		return extWEBP
	default:
		return ""
	}
}
