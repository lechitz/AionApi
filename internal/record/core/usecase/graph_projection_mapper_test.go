package usecase

import (
	"testing"
	"time"

	categorydomain "github.com/lechitz/aion-api/internal/category/core/domain"
	recorddomain "github.com/lechitz/aion-api/internal/record/core/domain"
	tagdomain "github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/stretchr/testify/require"
)

func TestBuildGraphProjection_BuildsStableNodesAndEdges(t *testing.T) {
	description := "Caminhada no fim do dia"
	status := "done"
	source := "manual"
	action := "Continue repetindo a rotina."
	generatedAt := time.Date(2026, 3, 11, 2, 0, 0, 0, time.UTC)

	projection := BuildGraphProjection(GraphProjectionBuildInput{
		UserID:      999,
		GeneratedAt: generatedAt,
		Timezone:    "America/Sao_Paulo",
		Categories: []categorydomain.Category{
			{ID: 1, UserID: 999, Name: "Saude", Description: "Saúde física", Color: "#16a34a", Icon: "heart"},
		},
		Tags: []tagdomain.Tag{
			{ID: 10, UserID: 999, CategoryID: 1, Name: "Agua", Description: "Água", Icon: "drop", UsageCount: 4},
		},
		Records: []recorddomain.Record{
			{
				ID:          77,
				UserID:      999,
				TagID:       10,
				Description: &description,
				EventTime:   time.Date(2026, 3, 10, 3, 30, 0, 0, time.UTC), // 00:30 America/Sao_Paulo
				Status:      &status,
				Source:      &source,
			},
		},
		Insights: []recorddomain.InsightCard{
			{
				ID:                "consistency-window_7d",
				Type:              "consistency_trend",
				Title:             "Consistencia forte na janela",
				Summary:           "Voce esteve ativo em 6 de 7 dias da janela.",
				Status:            "positive",
				Window:            recorddomain.InsightWindow7D,
				Confidence:        82,
				MetricKeys:        []string{"records.count"},
				RecommendedAction: &action,
				GeneratedAt:       generatedAt,
			},
		},
		InsightSupportedRecordIDs: map[string][]uint64{
			"consistency-window_7d": {77},
		},
		InsightScopedTagIDs: map[string][]uint64{
			"consistency-window_7d": {10},
		},
		InsightScopedCategoryIDs: map[string][]uint64{
			"consistency-window_7d": {1},
		},
	})

	require.Equal(t, recorddomain.GraphProjectionVersionV1, projection.Version)
	require.Equal(t, uint64(999), projection.UserID)
	require.Equal(t, generatedAt.UTC(), projection.GeneratedAt)

	nodeByID := make(map[string]recorddomain.GraphNode, len(projection.Nodes))
	for _, node := range projection.Nodes {
		nodeByID[node.ID] = node
	}
	require.Contains(t, nodeByID, "user:999")
	require.Contains(t, nodeByID, "category:1")
	require.Contains(t, nodeByID, "tag:10")
	require.Contains(t, nodeByID, "record:77")
	require.Contains(t, nodeByID, "time_bucket:day:2026-03-10")
	require.Contains(t, nodeByID, "insight:consistency-window_7d")

	require.Equal(t, recorddomain.GraphNodeTypeTimeBucket, nodeByID["time_bucket:day:2026-03-10"].Type)
	require.Equal(t, "day", nodeByID["time_bucket:day:2026-03-10"].Metadata["granularity"])
	require.Equal(t, "2026-03-10", nodeByID["record:77"].Metadata["event_date_local"])
	require.Equal(t, "America/Sao_Paulo", nodeByID["record:77"].Metadata["timezone"])

	edgeByID := make(map[string]recorddomain.GraphEdge, len(projection.Edges))
	for _, edge := range projection.Edges {
		edgeByID[edge.ID] = edge
	}
	require.Contains(t, edgeByID, "user_created_record:user:999:record:77")
	require.Contains(t, edgeByID, "record_tagged_as:record:77:tag:10")
	require.Contains(t, edgeByID, "tag_belongs_to_category:tag:10:category:1")
	require.Contains(t, edgeByID, "record_occurred_in:record:77:time_bucket:day:2026-03-10")
	require.Contains(t, edgeByID, "insight_supported_by_record:insight:consistency-window_7d:record:77")
	require.Contains(t, edgeByID, "insight_scoped_to_tag:insight:consistency-window_7d:tag:10")
	require.Contains(t, edgeByID, "insight_scoped_to_category:insight:consistency-window_7d:category:1")
}
