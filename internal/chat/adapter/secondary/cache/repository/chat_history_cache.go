// Package repository provides cache repository implementations for chat history.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
	output "github.com/lechitz/AionApi/internal/platform/ports/output/cache"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	// TracerName for observability.
	TracerName = "aionapi.chat.cache.repository"

	// SpanGetHistory is the span name for getting chat history from cache.
	SpanGetHistory = "chat.cache.history.get"
	// SpanSetHistory is the span name for setting chat history in cache.
	SpanSetHistory = "chat.cache.history.set"

	// LogGetHistory is the log message for getting chat history from cache.
	LogGetHistory = "Getting chat history from cache"
	// LogHistoryRetrieved is the log message for successfully retrieved chat history.
	LogHistoryRetrieved = "Chat history retrieved from cache"
	// LogHistoryCacheMiss is the log message for cache miss.
	LogHistoryCacheMiss = "Chat history cache miss"
	// LogSetHistory is the log message for setting chat history in cache.
	LogSetHistory = "Setting chat history in cache"
	// LogHistorySet is the log message for successfully set chat history.
	LogHistorySet = "Chat history set in cache"

	// ErrGetHistory is the error message for getting chat history from cache.
	ErrGetHistory = "failed to get chat history from cache"
	// ErrSetHistory is the error message for setting chat history in cache.
	ErrSetHistory = "failed to set chat history in cache"
	// ErrUnmarshalHistory is the error message for unmarshaling chat history.
	ErrUnmarshalHistory = "failed to unmarshal chat history"
	// ErrMarshalHistory is the error message for marshaling chat history.
	ErrMarshalHistory = "failed to marshal chat history"

	// Cache key prefix for chat history.
	cacheKeyPrefix = "chat:history:"
	// Default TTL for chat history cache (30 minutes).
	defaultTTL = 30 * time.Minute
)

// ChatHistoryCacheRepository provides cache operations for chat history.
type ChatHistoryCacheRepository struct {
	cache  output.Cache
	logger logger.ContextLogger
}

// NewChatHistoryCache creates a new ChatHistoryCacheRepository.
func NewChatHistoryCache(cache output.Cache, log logger.ContextLogger) *ChatHistoryCacheRepository {
	return &ChatHistoryCacheRepository{
		cache:  cache,
		logger: log,
	}
}

// buildKey builds a cache key for a user's chat history.
func buildKey(userID uint64) string {
	return fmt.Sprintf("%s%d", cacheKeyPrefix, userID)
}

// GetLatest retrieves the N most recent chat history entries for a user from cache.
func (r *ChatHistoryCacheRepository) GetLatest(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	r.logger.InfowCtx(ctx, LogGetHistory,
		commonkeys.UserID, userID,
		"limit", limit,
	)

	key := buildKey(userID)
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrGetHistory)
		r.logger.ErrorwCtx(ctx, ErrGetHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, fmt.Errorf("%s: %w", ErrGetHistory, err)
	}

	// Cache miss - return empty slice
	if data == "" {
		span.SetAttributes(attribute.Bool("cache_hit", false))
		r.logger.InfowCtx(ctx, LogHistoryCacheMiss, commonkeys.UserID, userID)
		return []domain.ChatHistory{}, nil
	}

	// Unmarshal cached history
	var history []domain.ChatHistory
	if err := json.Unmarshal([]byte(data), &history); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrUnmarshalHistory)
		r.logger.ErrorwCtx(ctx, ErrUnmarshalHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, fmt.Errorf("%s: %w", ErrUnmarshalHistory, err)
	}

	// Apply limit and return most recent
	if len(history) > limit {
		history = history[:limit]
	}

	span.SetAttributes(
		attribute.Bool("cache_hit", true),
		attribute.Int("results_count", len(history)),
	)
	span.SetStatus(codes.Ok, LogHistoryRetrieved)

	r.logger.InfowCtx(ctx, LogHistoryRetrieved,
		commonkeys.UserID, userID,
		"count", len(history),
	)

	return history, nil
}

// Add adds a new chat history entry to cache (prepends to the list).
func (r *ChatHistoryCacheRepository) Add(ctx context.Context, userID uint64, history domain.ChatHistory) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSetHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	r.logger.InfowCtx(ctx, LogSetHistory, commonkeys.UserID, userID)

	key := buildKey(userID)

	// Get existing history
	existingHistory, err := r.GetLatest(ctx, userID, 20) // Keep max 20 in cache
	if err != nil {
		// If error, start with empty history
		existingHistory = []domain.ChatHistory{}
	}

	// Prepend new entry (most recent first)
	newHistory := append([]domain.ChatHistory{history}, existingHistory...)

	// Limit to 20 most recent
	if len(newHistory) > 20 {
		newHistory = newHistory[:20]
	}

	// Marshal to JSON
	data, err := json.Marshal(newHistory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMarshalHistory)
		r.logger.ErrorwCtx(ctx, ErrMarshalHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return fmt.Errorf("%s: %w", ErrMarshalHistory, err)
	}

	// Set in cache with TTL
	if err := r.cache.Set(ctx, key, string(data), defaultTTL); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrSetHistory)
		r.logger.ErrorwCtx(ctx, ErrSetHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return fmt.Errorf("%s: %w", ErrSetHistory, err)
	}

	span.SetAttributes(attribute.Int("history_size", len(newHistory)))
	span.SetStatus(codes.Ok, LogHistorySet)

	r.logger.InfowCtx(ctx, LogHistorySet,
		commonkeys.UserID, userID,
		"history_size", len(newHistory),
	)

	return nil
}

// Clear removes chat history cache for a user.
func (r *ChatHistoryCacheRepository) Clear(ctx context.Context, userID uint64) error {
	key := buildKey(userID)
	return r.cache.Del(ctx, key)
}
