package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete performs a soft delete operation on a tag, marking it as deleted for the given user.
func (s *Service) SoftDelete(ctx context.Context, tagID uint64, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDeleteTag, trace.WithAttributes(
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, SpanSoftDeleteTag),
	))
	defer span.End()

	if tagID == 0 || userID == 0 {
		span.SetStatus(codes.Error, FailedToSoftDeleteTag)
		return errors.New(FailedToSoftDeleteTag)
	}

	if err := s.TagRepository.SoftDelete(ctx, tagID, userID); err != nil {
		span.SetStatus(codes.Error, FailedToSoftDeleteTag)
		span.RecordError(err)
		return errors.New(FailedToSoftDeleteTag)
	}

	span.AddEvent("InvalidateCache")
	err := s.TagCache.DeleteTag(ctx, tagID, userID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to invalidate tag cache after soft delete",
			commonkeys.TagID, strconv.FormatUint(tagID, 10),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err,
		)
	}

	err = s.TagCache.DeleteTagList(ctx, userID)
	if err != nil {
		s.Logger.WarnwCtx(ctx, "failed to invalidate tag list cache after soft delete",
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			commonkeys.Error, err,
		)
	}
	// Note: We can't invalidate by category here as we don't have the categoryID
	// The cache will expire naturally or be invalidated when the category is modified

	span.SetStatus(codes.Ok, StatusSoftDeleted)
	s.Logger.InfowCtx(ctx, SuccessfullySoftDeletedTag, commonkeys.TagID, strconv.FormatUint(tagID, 10))
	return nil
}
