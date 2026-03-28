package handler

import (
	"net/http"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/user/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const maxAvatarBytes = int64(20 << 20) // 20MB

// UploadAvatar uploads and validates user avatar image.
//
// @Summary      Upload user avatar
// @Description  Uploads an avatar image (PNG/JPEG/WEBP) and returns a URL payload suitable for user create/update.
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        avatar  formData  file  true  "Avatar image file"
// @Success      200     {object}  dto.UploadAvatarResponse
// @Failure      400     {string}  string  "Invalid file or validation error"
// @Failure      500     {string}  string  "Internal server error"
// @Router       /user/avatar/upload [post].
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer(TracerUserHandler)
	ctx, span := tr.Start(r.Context(), SpanUploadAvatarHandler)
	defer span.End()

	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarBytes+1024)
	if err := r.ParseMultipartForm(maxAvatarBytes); err != nil {
		span.SetStatus(codes.Error, "invalid multipart form")
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, sharederrors.NewValidationError("avatar", "invalid multipart form"), h.Logger)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		span.SetStatus(codes.Error, "avatar file is required")
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, sharederrors.NewValidationError("avatar", "avatar file is required"), h.Logger)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			h.Logger.WarnwCtx(ctx, "failed to close avatar file", "error", closeErr)
		}
	}()

	contentType := header.Header.Get("Content-Type")
	span.SetAttributes(
		attribute.String("avatar.filename", header.Filename),
		attribute.String("avatar.content_type", contentType),
		attribute.Int64("avatar.size_bytes", header.Size),
	)

	span.AddEvent(EventUserServiceUploadAvatar)
	avatarURL, normalizedContentType, sizeBytes, uploadErr := h.UserService.UploadAvatar(ctx, input.UploadAvatarCommand{
		File:        file,
		Filename:    header.Filename,
		SizeBytes:   header.Size,
		ContentType: contentType,
	})
	if uploadErr != nil {
		span.SetStatus(codes.Error, uploadErr.Error())
		httpresponse.WriteDomainErrorSpan(ctx, w, span, uploadErr, ErrUploadAvatar, h.Logger)
		return
	}

	span.SetAttributes(attribute.Int("http.status_code", http.StatusOK))
	span.SetStatus(codes.Ok, StatusAvatarUploaded)
	span.AddEvent(EventAvatarUploadedSuccess, trace.WithAttributes(attribute.Int64("avatar.size_bytes", sizeBytes)))

	httpresponse.WriteSuccess(w, http.StatusOK, dto.UploadAvatarResponse{
		AvatarURL:   avatarURL,
		ContentType: normalizedContentType,
		SizeBytes:   sizeBytes,
	}, "avatar uploaded successfully")
}
