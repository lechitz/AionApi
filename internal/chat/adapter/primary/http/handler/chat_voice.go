package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ChatVoice processes an audio message from the user and returns the AI response.
//
// @Summary      Send voice chat message
// @Description  Sends an audio file to be transcribed and processed by the AI assistant. Requires authentication.
// @Tags         ChatText
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header    string  true   "Bearer token"
// @Param        audio          formData  file    true   "Audio file (webm, wav, mp3, max 10MB, max 60s)"
// @Param        language       formData  string  false  "Language code (pt, en, es) or auto-detect if empty"
// @Success      200            {object}  map[string]interface{}  "Voice chat response with transcription and AI response"
// @Failure      400            {string}  string                  "Invalid audio file or validation error"
// @Failure      401            {string}  string                  "Unauthorized - missing or invalid token"
// @Failure      413            {string}  string                  "Audio file too large"
// @Failure      500            {string}  string                  "Internal server error"
// @Failure      503            {string}  string                  "Service unavailable - AI service is down"
// @Router       /chat/audio [post]
// @Security     BearerAuth.
func (h *Handler) ChatVoice(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerChatHandler).
		Start(r.Context(), SpanChatVoice)
	defer span.End()

	// Extract user ID from context
	userIDValue := ctx.Value(ctxkeys.UserID)
	if userIDValue == nil {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrUserIDNotFound), h.Logger)
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		span.SetStatus(codes.Error, ErrInvalidUserIDType)
		h.Logger.ErrorwCtx(ctx, LogInvalidUserIDType, "value", userIDValue)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewAuthenticationError(ErrInvalidUserID), h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(tracingkeys.RequestIPKey, r.RemoteAddr),
	)

	// Parse multipart form (max 10MB)
	span.AddEvent(EventParseMultipartForm)
	if err := r.ParseMultipartForm(MaxAudioSize); err != nil {
		span.SetStatus(codes.Error, ErrFailedParseMultipartForm)
		h.Logger.ErrorwCtx(ctx, LogFailedParseMultipartForm, commonkeys.Error, err)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio, ErrInvalidMultipartForm), h.Logger)
		return
	}

	// Get audio file
	file, header, err := r.FormFile(FormFieldAudio)
	if err != nil {
		span.SetStatus(codes.Error, ErrMissingAudioFile)
		h.Logger.ErrorwCtx(ctx, LogMissingAudioFile, commonkeys.Error, err)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio, ErrAudioFileRequired), h.Logger)
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > MaxAudioSize {
		span.SetStatus(codes.Error, ErrAudioFileTooLarge)
		h.Logger.ErrorwCtx(ctx, LogAudioFileTooLarge,
			"size", header.Size, "max", MaxAudioSize)
		httpresponse.WriteDecodeErrorSpan(ctx, w, span,
			sharederrors.NewValidationError(FormFieldAudio,
				fmt.Sprintf("Audio file too large: %d bytes (max: %d)", header.Size, MaxAudioSize)), h.Logger)
		return
	}

	span.SetAttributes(
		attribute.String(AttrAudioFilename, header.Filename),
		attribute.Int64(AttrAudioSizeBytes, header.Size),
		attribute.String(AttrAudioContentType, header.Header.Get(HeaderContentType)),
	)

	// Get optional language parameter
	language := r.FormValue(FormFieldLanguage)
	if language != "" {
		span.SetAttributes(attribute.String(AttrAudioLanguage, language))
	}

	// Forward to aion-chat service
	span.AddEvent(EventForwardToAionChat)

	aionChatURL := h.Config.AionChat.BaseURL + PathProcessAudio

	// Create new multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add audio file
	part, err := writer.CreateFormFile(FormFieldAudio, header.Filename)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCreateFormFile)
		h.Logger.ErrorwCtx(ctx, LogFailedCreateFormFile, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return
	}

	if _, err := io.Copy(part, file); err != nil {
		span.SetStatus(codes.Error, ErrFailedCopyAudioData)
		h.Logger.ErrorwCtx(ctx, LogFailedCopyAudioData, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgFailedProcessAudio, h.Logger)
		return
	}

	// Add form fields
	writer.WriteField(FormFieldUserID, strconv.FormatUint(userID, 10))
	if language != "" {
		writer.WriteField(FormFieldLanguage, language)
	}

	writer.Close()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, HTTPMethodPost, aionChatURL, &buf)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCreateRequest)
		h.Logger.ErrorwCtx(ctx, LogFailedCreateRequest, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgInternalServerError, h.Logger)
		return
	}

	req.Header.Set(HeaderContentType, writer.FormDataContentType())

	// Send request to aion-chat
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedCallAionChat)
		h.Logger.ErrorwCtx(ctx, LogFailedCallAionChat, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgAIServiceUnavailable, h.Logger)
		return
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int(AttrAionChatStatusCode, resp.StatusCode))

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(codes.Error, ErrFailedReadResponse)
		h.Logger.ErrorwCtx(ctx, LogFailedReadResponse, commonkeys.Error, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, MsgInternalServerError, h.Logger)
		return
	}

	// Forward response status and body
	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, ErrAionChatReturnedError)
		h.Logger.ErrorwCtx(ctx, LogAionChatError,
			"status_code", resp.StatusCode, "response", string(responseBody))
		w.Header().Set(HeaderContentType, ContentTypeJSON)
		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)
		return
	}

	// Success
	span.SetStatus(codes.Ok, StatusVoiceChatSuccess)
	h.Logger.InfowCtx(ctx, LogVoiceChatSuccess,
		commonkeys.UserID, userID,
		"audio_size", header.Size,
		"status_code", resp.StatusCode)

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
