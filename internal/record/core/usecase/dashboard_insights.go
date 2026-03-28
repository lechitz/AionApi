package usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// InsightFeed returns deterministic, explainable insights for a given analysis window.
func (s *Service) InsightFeed(ctx context.Context, userID uint64, query input.InsightFeedQuery) ([]domain.InsightCard, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanInsightFeed)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanInsightFeed),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrWindow, query.Window),
		attribute.Int(AttrLimit, query.Limit),
		attribute.Int(AttrTagIDsCount, len(query.TagIDs)),
		attribute.String(AttrTimezone, query.Timezone),
	)
	if userID == 0 {
		span.RecordError(ErrUserIDIsRequired)
		span.SetStatus(codes.Error, UserIDIsRequired)
		s.Logger.ErrorwCtx(ctx, UserIDIsRequired)
		return nil, ErrUserIDIsRequired
	}
	if query.CategoryID != nil {
		span.SetAttributes(attribute.String(commonkeys.CategoryID, strconv.FormatUint(*query.CategoryID, 10)))
	}

	window := normalizeInsightWindow(query.Window)
	limit := normalizeInsightLimit(query.Limit)
	loc, _ := resolveInsightLocation(query.Timezone)
	targetDate := normalizeInsightDate(query.Date, loc)

	startUTC, endUTC, windowDays := insightRangeForWindow(targetDate, loc, window)
	prevStartUTC, prevEndUTC := insightPreviousRange(startUTC, windowDays)

	records, err := s.RecordRepository.ListAllBetween(ctx, userID, startUTC, endUTC, DefaultDashboardLimit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeInsightFeed)
		s.Logger.ErrorwCtx(ctx, ErrComputeInsightFeed, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}
	prevRecords, err := s.RecordRepository.ListAllBetween(ctx, userID, prevStartUTC, prevEndUTC, DefaultDashboardLimit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeInsightFeed)
		s.Logger.ErrorwCtx(ctx, ErrComputeInsightFeed, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}
	span.AddEvent(EventRepositoryMetricDefinitions)
	defs, err := s.RecordRepository.ListMetricDefinitions(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeInsightFeed)
		s.Logger.ErrorwCtx(ctx, ErrComputeInsightFeed, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}
	span.AddEvent(EventRepositoryTags)
	tags, err := s.TagRepository.GetAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeInsightFeed)
		s.Logger.ErrorwCtx(ctx, ErrComputeInsightFeed, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return nil, err
	}

	records = filterInsightRecordsByScope(records, query.CategoryID, query.TagIDs, tags)
	prevRecords = filterInsightRecordsByScope(prevRecords, query.CategoryID, query.TagIDs, tags)

	now := time.Now().UTC()
	insights := make([]domain.InsightCard, 0, 5)

	if item := buildConsistencyTrendInsight(records, window, now, loc, windowDays); item != nil {
		insights = append(insights, *item)
	}
	if item := buildStreakRiskInsight(records, window, now, loc, targetDate); item != nil {
		insights = append(insights, *item)
	}
	if item := buildActivityGapInsight(records, window, now, loc, targetDate); item != nil {
		insights = append(insights, *item)
	}
	if item := buildCategoryConcentrationInsight(records, defs, window, now, query.CategoryID, query.TagIDs); item != nil {
		insights = append(insights, *item)
	}
	if item := buildRecentChangeInsight(records, prevRecords, window, now, windowDays); item != nil {
		insights = append(insights, *item)
	}

	sort.SliceStable(insights, func(i, j int) bool {
		if insights[i].Confidence == insights[j].Confidence {
			return insights[i].Type < insights[j].Type
		}
		return insights[i].Confidence > insights[j].Confidence
	})

	if len(insights) > limit {
		insights = insights[:limit]
	}
	span.AddEvent(EventSuccess)
	span.SetAttributes(attribute.Int(AttrResultsCount, len(insights)))
	span.SetStatus(codes.Ok, StatusFetched)
	s.Logger.InfowCtx(ctx, LogInsightFeedComputedSuccessfully,
		commonkeys.UserID, userID,
		AttrResultsCount, len(insights),
		AttrWindow, string(window),
	)
	return insights, nil
}

// AnalyticsSeries returns a compact time series for a supported key/window pair.
func (s *Service) AnalyticsSeries(ctx context.Context, userID uint64, query input.AnalyticsSeriesQuery) (domain.AnalyticsSeriesResult, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanAnalyticsSeries)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanAnalyticsSeries),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(AttrSeriesKey, query.SeriesKey),
		attribute.String(AttrWindow, query.Window),
		attribute.Int(AttrTagIDsCount, len(query.TagIDs)),
		attribute.String(AttrTimezone, query.Timezone),
	)
	if userID == 0 {
		span.RecordError(ErrUserIDIsRequired)
		span.SetStatus(codes.Error, UserIDIsRequired)
		s.Logger.ErrorwCtx(ctx, UserIDIsRequired)
		return domain.AnalyticsSeriesResult{}, ErrUserIDIsRequired
	}
	if query.CategoryID != nil {
		span.SetAttributes(attribute.String(commonkeys.CategoryID, strconv.FormatUint(*query.CategoryID, 10)))
	}

	window := normalizeInsightWindow(query.Window)
	loc, _ := resolveInsightLocation(query.Timezone)
	targetDate := normalizeInsightDate(query.Date, loc)
	startUTC, _, windowDays := insightRangeForWindow(targetDate, loc, window)
	endUTC := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()

	records, err := s.RecordRepository.ListAllBetween(ctx, userID, startUTC, endUTC, DefaultDashboardLimit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeAnalyticsSeries)
		s.Logger.ErrorwCtx(ctx, ErrComputeAnalyticsSeries, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return domain.AnalyticsSeriesResult{}, err
	}
	span.AddEvent(EventRepositoryMetricDefinitions)
	defs, err := s.RecordRepository.ListMetricDefinitions(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeAnalyticsSeries)
		s.Logger.ErrorwCtx(ctx, ErrComputeAnalyticsSeries, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return domain.AnalyticsSeriesResult{}, err
	}
	span.AddEvent(EventRepositoryTags)
	tags, err := s.TagRepository.GetAll(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrComputeAnalyticsSeries)
		s.Logger.ErrorwCtx(ctx, ErrComputeAnalyticsSeries, commonkeys.Error, err.Error(), commonkeys.UserID, userID)
		return domain.AnalyticsSeriesResult{}, err
	}
	records = filterInsightRecordsByScope(records, query.CategoryID, query.TagIDs, tags)

	seriesKey := strings.TrimSpace(query.SeriesKey)
	points := make([]domain.AnalyticsPoint, 0, windowDays)
	for i := range windowDays {
		dayLocal := targetDate.AddDate(0, 0, -(windowDays - 1 - i))
		dayStartUTC := time.Date(dayLocal.Year(), dayLocal.Month(), dayLocal.Day(), 0, 0, 0, 0, loc).UTC()
		dayEndUTC := time.Date(dayLocal.Year(), dayLocal.Month(), dayLocal.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()
		dayRecords := filterRecordsBetween(records, dayStartUTC, dayEndUTC)
		value := analyticsValueForSeries(dayRecords, defs, seriesKey)
		label := dayLocal.Format("2006-01-02")
		points = append(points, domain.AnalyticsPoint{
			Timestamp: dayStartUTC,
			Value:     &value,
			Label:     &label,
		})
	}

	summary := fmt.Sprintf("%s across %d days", seriesKey, windowDays)
	result := domain.AnalyticsSeriesResult{
		SeriesKey: seriesKey,
		Window:    window,
		Points:    points,
		Summary:   &summary,
	}
	span.AddEvent(EventSuccess)
	span.SetAttributes(attribute.Int(AttrResultsCount, len(points)))
	span.SetStatus(codes.Ok, StatusStatsComputed)
	s.Logger.InfowCtx(ctx, LogAnalyticsSeriesComputedSuccessfully,
		commonkeys.UserID, userID,
		AttrSeriesKey, seriesKey,
		AttrResultsCount, len(points),
	)
	return result, nil
}

func normalizeInsightWindow(raw string) domain.InsightWindow {
	switch strings.TrimSpace(strings.ToUpper(raw)) {
	case string(domain.InsightWindow30D):
		return domain.InsightWindow30D
	case string(domain.InsightWindow90D):
		return domain.InsightWindow90D
	default:
		return domain.InsightWindow7D
	}
}

func normalizeInsightLimit(limit int) int {
	if limit <= 0 {
		return DefaultInsightFeedLimit
	}
	if limit > MaxInsightFeedLimit {
		return MaxInsightFeedLimit
	}
	return limit
}

func resolveInsightLocation(timezone string) (*time.Location, string) {
	tzName := strings.TrimSpace(timezone)
	if tzName == "" {
		tzName = DefaultTimezone
	}
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.UTC, "UTC"
	}
	return loc, tzName
}

func normalizeInsightDate(date time.Time, loc *time.Location) time.Time {
	return normalizeDashboardDate(date, loc)
}

func insightRangeForWindow(targetDate time.Time, loc *time.Location, window domain.InsightWindow) (time.Time, time.Time, int) {
	days := insightWindowDays(window)
	startLocal := targetDate.AddDate(0, 0, -(days - 1))
	endLocal := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), loc)
	return startLocal.UTC(), endLocal.UTC(), days
}

func insightPreviousRange(currentStartUTC time.Time, days int) (time.Time, time.Time) {
	prevEnd := currentStartUTC.Add(-time.Nanosecond)
	prevStart := currentStartUTC.AddDate(0, 0, -days)
	return prevStart, prevEnd
}

func insightWindowDays(window domain.InsightWindow) int {
	switch window {
	case domain.InsightWindow7D:
		return 7
	case domain.InsightWindow30D:
		return 30
	case domain.InsightWindow90D:
		return 90
	}
	return 7
}

func buildConsistencyTrendInsight(records []domain.Record, window domain.InsightWindow, now time.Time, loc *time.Location, windowDays int) *domain.InsightCard {
	activeDays := len(activeLocalDays(records, loc))
	if activeDays == 0 {
		return nil
	}

	ratio := float64(activeDays) / float64(windowDays)
	status := "neutral"
	confidence := 60
	title := "Consistencia em desenvolvimento"
	summary := fmt.Sprintf("Voce registrou atividade em %d de %d dias.", activeDays, windowDays)
	action := "Mantenha uma pequena acao diaria para preservar a continuidade."

	if ratio >= 0.7 {
		status = "positive"
		confidence = 82
		title = "Consistencia forte na janela"
		summary = fmt.Sprintf("Voce esteve ativo em %d de %d dias da janela.", activeDays, windowDays)
		action = "Continue repetindo a rotina que manteve a frequencia alta."
	} else if ratio <= 0.3 {
		status = "warning"
		confidence = 74
		title = "Consistencia abaixo do ideal"
		summary = fmt.Sprintf("A atividade apareceu em apenas %d de %d dias.", activeDays, windowDays)
		action = "Reduza a friccao e registre pelo menos uma acao pequena por dia."
	}

	return &domain.InsightCard{
		ID:                fmt.Sprintf("consistency-%s", strings.ToLower(string(window))),
		Type:              "consistency_trend",
		Title:             title,
		Summary:           summary,
		Status:            status,
		Window:            window,
		Confidence:        confidence,
		MetricKeys:        []string{"records.count"},
		RecommendedAction: strPtr(action),
		Evidence: []domain.InsightEvidence{
			{Label: "dias ativos", Value: strconv.Itoa(activeDays), Kind: "count"},
			{Label: "janela", Value: fmt.Sprintf("%d dias", windowDays), Kind: "window"},
		},
		GeneratedAt: now,
	}
}

func buildStreakRiskInsight(records []domain.Record, window domain.InsightWindow, now time.Time, loc *time.Location, targetDate time.Time) *domain.InsightCard {
	last, ok := latestRecord(records)
	if !ok {
		return nil
	}

	idleDays := int(targetDate.Sub(last.EventTime.In(loc)).Hours() / 24)
	if idleDays < 2 {
		return nil
	}

	action := "Retome hoje com uma acao pequena para nao alongar a pausa."
	return &domain.InsightCard{
		ID:                fmt.Sprintf("streak-risk-%s", strings.ToLower(string(window))),
		Type:              "streak_risk",
		Title:             "Risco de quebra de ritmo",
		Summary:           fmt.Sprintf("Seu ultimo registro ocorreu ha %d dias.", idleDays),
		Status:            "warning",
		Window:            window,
		Confidence:        76,
		MetricKeys:        []string{"records.count"},
		RecommendedAction: strPtr(action),
		Evidence: []domain.InsightEvidence{
			{Label: "ultimo registro", Value: last.EventTime.In(loc).Format("2006-01-02"), Kind: "date"},
			{Label: "dias sem atividade", Value: strconv.Itoa(idleDays), Kind: "count"},
		},
		GeneratedAt: now,
	}
}

func buildActivityGapInsight(records []domain.Record, window domain.InsightWindow, now time.Time, loc *time.Location, targetDate time.Time) *domain.InsightCard {
	if len(records) > 0 {
		return nil
	}

	action := "Use o chat ou quick add para registrar o primeiro evento da janela."
	return &domain.InsightCard{
		ID:                fmt.Sprintf("activity-gap-%s", strings.ToLower(string(window))),
		Type:              "activity_gap",
		Title:             "Sem atividade na janela analisada",
		Summary:           "Nao houve registros no periodo solicitado.",
		Status:            "warning",
		Window:            window,
		Confidence:        90,
		MetricKeys:        []string{"records.count"},
		RecommendedAction: strPtr(action),
		Evidence: []domain.InsightEvidence{
			{Label: "janela final", Value: targetDate.In(loc).Format("2006-01-02"), Kind: "date"},
			{Label: "registros", Value: "0", Kind: "count"},
		},
		GeneratedAt: now,
	}
}

func buildCategoryConcentrationInsight(
	records []domain.Record,
	defs []domain.MetricDefinition,
	window domain.InsightWindow,
	now time.Time,
	categoryID *uint64,
	tagIDs []uint64,
) *domain.InsightCard {
	if categoryID != nil || len(tagIDs) > 0 {
		return nil
	}
	if len(records) == 0 || len(defs) == 0 {
		return nil
	}

	counts := make(map[string]int)
	tagToName := make(map[uint64]string)
	for _, def := range defs {
		for _, tagID := range def.TagIDs {
			tagToName[tagID] = def.DisplayName
		}
		tagToName[def.TagID] = def.DisplayName
	}

	for _, rec := range records {
		name := tagToName[rec.TagID]
		if name == "" {
			name = "Sem metrica"
		}
		counts[name]++
	}

	var topName string
	topCount := 0
	for name, count := range counts {
		if count > topCount {
			topName = name
			topCount = count
		}
	}
	if topCount == 0 {
		return nil
	}

	share := float64(topCount) / float64(len(records))
	if share < 0.6 {
		return nil
	}

	action := "Avalie se essa concentracao foi intencional ou se vale reequilibrar a rotina."
	return &domain.InsightCard{
		ID:                fmt.Sprintf("category-concentration-%s", strings.ToLower(string(window))),
		Type:              "category_concentration",
		Title:             "Concentracao alta em uma frente",
		Summary:           fmt.Sprintf("%s representa %.0f%% dos registros da janela.", topName, share*100),
		Status:            "neutral",
		Window:            window,
		Confidence:        71,
		MetricKeys:        []string{"records.count"},
		RecommendedAction: strPtr(action),
		Evidence: []domain.InsightEvidence{
			{Label: "frente dominante", Value: topName, Kind: "label"},
			{Label: "participacao", Value: fmt.Sprintf("%.0f%%", share*100), Kind: "percentage"},
		},
		GeneratedAt: now,
	}
}

func buildRecentChangeInsight(records []domain.Record, prevRecords []domain.Record, window domain.InsightWindow, now time.Time, windowDays int) *domain.InsightCard {
	if len(prevRecords) == 0 {
		return nil
	}

	curr := float64(len(records))
	prev := float64(len(prevRecords))
	if prev <= 0 {
		return nil
	}

	deltaPct := ((curr - prev) / prev) * 100
	if math.Abs(deltaPct) < 20 {
		return nil
	}

	status := "positive"
	title := "Ritmo acima da janela anterior"
	action := "Verifique o que mudou recentemente para repetir o padrao."
	if deltaPct < 0 {
		status = "warning"
		title = "Ritmo abaixo da janela anterior"
		action = "Compare esta janela com a anterior e remova o principal ponto de friccao."
	}

	return &domain.InsightCard{
		ID:                fmt.Sprintf("recent-change-%s", strings.ToLower(string(window))),
		Type:              "recent_change",
		Title:             title,
		Summary:           fmt.Sprintf("A atividade variou %.0f%% em relacao aos %d dias anteriores.", deltaPct, windowDays),
		Status:            status,
		Window:            window,
		Confidence:        78,
		MetricKeys:        []string{"records.count"},
		RecommendedAction: strPtr(action),
		Evidence: []domain.InsightEvidence{
			{Label: "janela atual", Value: fmt.Sprintf("%.0f", curr), Kind: "count"},
			{Label: "janela anterior", Value: fmt.Sprintf("%.0f", prev), Kind: "count"},
		},
		GeneratedAt: now,
	}
}

func activeLocalDays(records []domain.Record, loc *time.Location) map[string]struct{} {
	out := make(map[string]struct{}, len(records))
	for _, rec := range records {
		key := rec.EventTime.In(loc).Format("2006-01-02")
		out[key] = struct{}{}
	}
	return out
}

func latestRecord(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	latest := records[0]
	for _, rec := range records[1:] {
		if rec.EventTime.After(latest.EventTime) {
			latest = rec
		}
	}
	return latest, true
}

func filterRecordsBetween(records []domain.Record, startUTC, endUTC time.Time) []domain.Record {
	out := make([]domain.Record, 0)
	for _, rec := range records {
		if rec.EventTime.Before(startUTC) || rec.EventTime.After(endUTC) {
			continue
		}
		out = append(out, rec)
	}
	return out
}

func analyticsValueForSeries(records []domain.Record, defs []domain.MetricDefinition, seriesKey string) float64 {
	if strings.TrimSpace(seriesKey) == "" || seriesKey == "records.count" {
		return float64(len(records))
	}

	var def *domain.MetricDefinition
	for i := range defs {
		if defs[i].MetricKey == seriesKey {
			def = &defs[i]
			break
		}
	}
	if def == nil {
		return float64(len(records))
	}
	return computeMetricValue(records, *def)
}

func filterInsightRecordsByScope(records []domain.Record, categoryID *uint64, tagIDs []uint64, tags []tagdomain.Tag) []domain.Record {
	if categoryID == nil && len(tagIDs) == 0 {
		return records
	}

	tagFilter := make(map[uint64]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID == 0 {
			continue
		}
		tagFilter[tagID] = struct{}{}
	}

	tagToCategory := make(map[uint64]uint64, len(tags))
	for _, tag := range tags {
		tagToCategory[tag.ID] = tag.CategoryID
	}

	out := make([]domain.Record, 0, len(records))
	for _, rec := range records {
		if categoryID != nil {
			if tagToCategory[rec.TagID] != *categoryID {
				continue
			}
		}
		if len(tagFilter) > 0 {
			if _, ok := tagFilter[rec.TagID]; !ok {
				continue
			}
		}
		out = append(out, rec)
	}
	return out
}

func strPtr(v string) *string {
	return &v
}
