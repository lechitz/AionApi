package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetUserStats retrieves aggregated statistics for a user.
func (s *Service) GetUserStats(ctx context.Context, userID uint64) (domain.UserStats, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetUserStats)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	s.logger.InfowCtx(ctx, InfoGettingUserStats, commonkeys.UserID, userID)

	stats, err := s.userRepository.GetUserStats(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToGetUserStats)
		s.logger.ErrorwCtx(ctx, ErrorToGetUserStats,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return domain.UserStats{}, err
	}

	span.SetAttributes(
		attribute.Int("total_records", stats.TotalRecords),
		attribute.Int("records_this_month", stats.RecordsThisMonth),
	)
	span.SetStatus(codes.Ok, SuccessUserStatsRetrieved)

	s.logger.InfowCtx(ctx, SuccessUserStatsRetrieved,
		commonkeys.UserID, userID,
		"total_records", stats.TotalRecords,
		"records_this_month", stats.RecordsThisMonth,
	)

	return stats, nil
}
