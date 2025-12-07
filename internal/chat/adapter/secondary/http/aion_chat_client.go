// Package http provides the HTTP client adapter for communicating with Aion-Chat service.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lechitz/AionApi/internal/chat/adapter/primary/http/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SendMessage sends a chat message to the Aion-Chat service.
func (c *AionChatClient) SendMessage(ctx context.Context, req *dto.InternalChatRequest) (*dto.InternalChatResponse, error) {
	tr := otel.Tracer(TracerAionChatClient)
	ctx, span := tr.Start(ctx, SpanSendMessage)
	defer span.End()

	url := fmt.Sprintf("%s%s", c.baseURL, PathProcess)

	span.SetAttributes(
		attribute.String(AttrHTTPURL, url),
		attribute.String(AttrHTTPMethod, http.MethodPost),
		attribute.Int64(AttrUserID, int64(req.UserID)),
	)

	jsonData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedMarshal)
		c.logger.ErrorwCtx(ctx, ErrFailedMarshal, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedMarshal, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedCreateRequest)
		c.logger.ErrorwCtx(ctx, ErrFailedCreateRequest, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedCreateRequest, err)
	}

	httpReq.Header.Set(HeaderContentType, ContentTypeJSON)
	httpReq.Header.Set(HeaderAccept, ContentTypeJSON)

	c.logger.InfowCtx(ctx, MsgCallingAionChatService, commonkeys.URL, url, AttrUserID, req.UserID)
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrHTTPRequestFailed)
		c.logger.ErrorwCtx(ctx, ErrHTTPRequestFailed, commonkeys.Error, err.Error(), "url", url)
		return nil, fmt.Errorf("%s: %w", ErrAionChatRequestFailed, err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	span.SetAttributes(attribute.Int(AttrHTTPStatusCode, httpResp.StatusCode))

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedReadResponse)
		c.logger.ErrorwCtx(ctx, ErrFailedReadResponse, commonkeys.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", ErrFailedReadResponse, err)
	}

	if httpResp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, ErrAionChatNonOK)
		c.logger.ErrorwCtx(ctx, ErrAionChatNonOK, "status_code", httpResp.StatusCode, "body", string(body))
		return nil, fmt.Errorf("%s: status %d: %s", ErrAionChatNonOK, httpResp.StatusCode, string(body))
	}

	var resp dto.InternalChatResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrFailedUnmarshal)
		c.logger.ErrorwCtx(ctx, ErrFailedUnmarshal, commonkeys.Error, err.Error(), "body", string(body))
		return nil, fmt.Errorf("%s: %w", ErrFailedUnmarshal, err)
	}

	span.SetAttributes(
		attribute.Int(AttrTokensUsed, resp.TokensUsed),
		attribute.Int(AttrResponseLength, len(resp.Response)),
	)
	span.SetStatus(codes.Ok, StatusMessageSent)

	c.logger.InfowCtx(ctx, MsgAionChatResponseReceived, AttrUserID, req.UserID, "response_length", len(resp.Response))

	return &resp, nil
}
