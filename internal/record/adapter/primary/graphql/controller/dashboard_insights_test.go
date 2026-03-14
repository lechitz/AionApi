package controller_test

import (
	"context"
	"testing"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsightFeed_Success_MapsEvidenceAndScopeQuery(t *testing.T) {
	var captured input.InsightFeedQuery
	generatedAt := time.Date(2026, 3, 11, 1, 45, 0, 0, time.UTC)
	svc := &recordServiceStub{
		insightFeedFn: func(_ context.Context, userID uint64, query input.InsightFeedQuery) ([]domain.InsightCard, error) {
			require.Equal(t, uint64(999), userID)
			captured = query
			return []domain.InsightCard{
				{
					ID:         "activity-gap-window_7d",
					Type:       "activity_gap",
					Title:      "Sem atividade na janela analisada",
					Summary:    "Nao houve registros no periodo solicitado.",
					Status:     "warning",
					Window:     domain.InsightWindow7D,
					Confidence: 90,
					MetricKeys: []string{"records.count"},
					RecommendedAction: func() *string {
						s := "Use o chat ou quick add para registrar o primeiro evento da janela."
						return &s
					}(),
					Evidence: []domain.InsightEvidence{
						{Label: "janela final", Value: "2026-03-10", Kind: "date"},
						{Label: "registros", Value: "0", Kind: "count"},
					},
					GeneratedAt: generatedAt,
				},
			}, nil
		},
	}

	h, ctrl := newRecordController(t, svc)
	defer ctrl.Finish()

	limit := int32(3)
	date := "2026-03-10"
	timezone := "America/Sao_Paulo"
	categoryID := "12"
	out, err := h.InsightFeed(t.Context(), 999, gmodel.InsightWindowWindow7d, &limit, &date, &timezone, &categoryID, []string{"4", "5"})

	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "activity_gap", out[0].Type)
	assert.Equal(t, "Sem atividade na janela analisada", out[0].Title)
	assert.Equal(t, int32(90), out[0].Confidence)
	require.Len(t, out[0].Evidence, 2)
	assert.Equal(t, "janela final", out[0].Evidence[0].Label)
	assert.Equal(t, "2026-03-10", out[0].Evidence[0].Value)
	assert.Equal(t, "count", out[0].Evidence[1].Kind)
	assert.Equal(t, generatedAt.Format(time.RFC3339), out[0].GeneratedAt)

	assert.Equal(t, string(domain.InsightWindow7D), captured.Window)
	assert.Equal(t, 3, captured.Limit)
	assert.Equal(t, "America/Sao_Paulo", captured.Timezone)
	require.NotNil(t, captured.CategoryID)
	assert.Equal(t, uint64(12), *captured.CategoryID)
	assert.Equal(t, []uint64{4, 5}, captured.TagIDs)
	assert.Equal(t, time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), captured.Date)
}
