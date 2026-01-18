package repository

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/chat/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/chat/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetLatest retrieves the N most recent chat history entries for a given user.
func (r *ChatHistoryRepository) GetLatest(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetLatestRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OpGetLatest),
		attribute.Int("limit", limit),
	))
	defer span.End()

	var chatHistoryDB []model.ChatHistoryDB

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&chatHistoryDB).Error()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetLatest)
		r.logger.ErrorwCtx(ctx, ErrGetLatestChatMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			"limit", limit,
		)
		return nil, err
	}

	histories := mapper.ChatHistoriesFromDB(chatHistoryDB)

	span.SetAttributes(attribute.Int("results_count", len(histories)))
	span.SetStatus(codes.Ok, StatusRetrievedLatest)

	r.logger.InfowCtx(ctx, StatusRetrievedLatest,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"count", len(histories),
	)

	return histories, nil
}

// GetByUserID retrieves chat history entries for a given user with pagination support.
func (r *ChatHistoryRepository) GetByUserID(ctx context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByUserIDRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OpGetByUserID),
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	))
	defer span.End()

	var chatHistoryDB []model.ChatHistoryDB

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&chatHistoryDB).Error()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpGetByUserID)
		r.logger.ErrorwCtx(ctx, ErrGetChatByUserIDMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			"limit", limit,
			"offset", offset,
		)
		return nil, err
	}

	histories := mapper.ChatHistoriesFromDB(chatHistoryDB)

	span.SetAttributes(attribute.Int("results_count", len(histories)))
	span.SetStatus(codes.Ok, StatusRetrievedByUserID)

	r.logger.InfowCtx(ctx, StatusRetrievedByUserID,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"count", len(histories),
		"limit", limit,
		"offset", offset,
	)

	return histories, nil
}
