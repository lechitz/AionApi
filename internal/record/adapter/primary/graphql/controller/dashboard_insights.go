package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (c *controller) InsightFeed(
	ctx context.Context,
	userID uint64,
	window model.InsightWindow,
	limit *int32,
	date *string,
	timezone *string,
	categoryID *string,
	tagIDs []string,
) ([]*model.InsightCard, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanInsightFeed)
	defer span.End()

	lim := 8
	if limit != nil && *limit > 0 {
		lim = int(*limit)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanInsightFeed),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrWindow, string(window)),
		attribute.Int(AttrLimit, lim),
		attribute.Int(AttrTagIDsCount, len(tagIDs)),
		attribute.String(AttrTimezone, stringOrEmpty(timezone)),
	)

	targetDate, err := parseDateOrDefault(stringOrEmpty(date))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgInvalidDateFormat)
		c.Logger.ErrorwCtx(ctx, MsgInvalidDateFormat, AttrDate, stringOrEmpty(date), commonkeys.Error, err.Error())
		return nil, err
	}
	if categoryID != nil {
		span.SetAttributes(attribute.String(commonkeys.CategoryID, *categoryID))
	}

	items, err := c.RecordService.InsightFeed(ctx, userID, input.InsightFeedQuery{
		Window:     string(window),
		Limit:      lim,
		Date:       targetDate,
		Timezone:   stringOrEmpty(timezone),
		CategoryID: parseOptionalID(categoryID),
		TagIDs:     parseIDs(tagIDs),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgInsightFeedError)
		c.Logger.ErrorwCtx(ctx, MsgInsightFeedError, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return nil, err
	}

	out := make([]*model.InsightCard, 0, len(items))
	for _, item := range items {
		out = append(out, toGraphQLInsightCard(item))
	}
	span.SetAttributes(attribute.Int(AttrResultsCount, len(out)))
	span.SetStatus(codes.Ok, StatusFetched)
	c.Logger.InfowCtx(ctx, MsgInsightFeedFetched,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		AttrResultsCount, len(out),
		AttrWindow, string(window),
	)
	return out, nil
}

func (c *controller) AnalyticsSeries(
	ctx context.Context,
	userID uint64,
	seriesKey string,
	window model.InsightWindow,
	date *string,
	timezone *string,
	categoryID *string,
	tagIDs []string,
) (*model.AnalyticsSeriesResult, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanAnalyticsSeries)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanAnalyticsSeries),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrSeriesKey, seriesKey),
		attribute.String(AttrWindow, string(window)),
		attribute.Int(AttrTagIDsCount, len(tagIDs)),
		attribute.String(AttrTimezone, stringOrEmpty(timezone)),
	)

	targetDate, err := parseDateOrDefault(stringOrEmpty(date))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgInvalidDateFormat)
		c.Logger.ErrorwCtx(ctx, MsgInvalidDateFormat, AttrDate, stringOrEmpty(date), commonkeys.Error, err.Error())
		return nil, err
	}
	if categoryID != nil {
		span.SetAttributes(attribute.String(commonkeys.CategoryID, *categoryID))
	}

	out, err := c.RecordService.AnalyticsSeries(ctx, userID, input.AnalyticsSeriesQuery{
		SeriesKey:  seriesKey,
		Window:     string(window),
		Date:       targetDate,
		Timezone:   stringOrEmpty(timezone),
		CategoryID: parseOptionalID(categoryID),
		TagIDs:     parseIDs(tagIDs),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgAnalyticsSeriesError)
		c.Logger.ErrorwCtx(ctx, MsgAnalyticsSeriesError, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10), AttrSeriesKey, seriesKey)
		return nil, err
	}
	span.SetAttributes(attribute.Int(AttrResultsCount, len(out.Points)))
	span.SetStatus(codes.Ok, StatusStatsComputed)
	c.Logger.InfowCtx(ctx, MsgAnalyticsSeriesComputed,
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		AttrSeriesKey, seriesKey,
		AttrResultsCount, len(out.Points),
	)
	return toGraphQLAnalyticsSeriesResult(out), nil
}

func toGraphQLInsightCard(in domain.InsightCard) *model.InsightCard {
	evidence := make([]*model.InsightEvidence, 0, len(in.Evidence))
	for _, item := range in.Evidence {
		evidence = append(evidence, &model.InsightEvidence{
			Label: item.Label,
			Value: item.Value,
			Kind:  item.Kind,
		})
	}

	return &model.InsightCard{
		ID:                in.ID,
		Type:              in.Type,
		Title:             in.Title,
		Summary:           in.Summary,
		Status:            in.Status,
		Window:            toGraphQLInsightWindow(in.Window),
		Confidence:        safeIntToInt32(in.Confidence),
		MetricKeys:        in.MetricKeys,
		RecommendedAction: in.RecommendedAction,
		Evidence:          evidence,
		GeneratedAt:       in.GeneratedAt.Format(timeLayout),
	}
}

func toGraphQLAnalyticsSeriesResult(in domain.AnalyticsSeriesResult) *model.AnalyticsSeriesResult {
	points := make([]*model.AnalyticsPoint, 0, len(in.Points))
	for _, point := range in.Points {
		points = append(points, &model.AnalyticsPoint{
			Timestamp: point.Timestamp.Format(timeLayout),
			Value:     point.Value,
			Label:     point.Label,
		})
	}

	return &model.AnalyticsSeriesResult{
		SeriesKey: in.SeriesKey,
		Window:    toGraphQLInsightWindow(in.Window),
		Points:    points,
		Summary:   in.Summary,
	}
}

func toGraphQLInsightWindow(v domain.InsightWindow) model.InsightWindow {
	switch v {
	case domain.InsightWindow7D:
		return model.InsightWindowWindow7d
	case domain.InsightWindow30D:
		return model.InsightWindowWindow30d
	case domain.InsightWindow90D:
		return model.InsightWindowWindow90d
	}
	return model.InsightWindowWindow7d
}

func stringOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func parseOptionalID(value *string) *uint64 {
	if value == nil || *value == "" {
		return nil
	}
	parsed, err := strconv.ParseUint(*value, 10, 64)
	if err != nil || parsed == 0 {
		return nil
	}
	return &parsed
}
