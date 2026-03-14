package domain

import "testing"

func TestGraphProjectionV1Constants(t *testing.T) {
	if GraphProjectionVersionV1 != "graph-projection-v1" {
		t.Fatalf("unexpected graph projection version: %s", GraphProjectionVersionV1)
	}

	if GraphNodeTypeInsight != "insight" {
		t.Fatalf("unexpected insight node type: %s", GraphNodeTypeInsight)
	}

	if GraphEdgeTypeInsightSupportedByRecord != "insight_supported_by_record" {
		t.Fatalf("unexpected insight-supported edge type: %s", GraphEdgeTypeInsightSupportedByRecord)
	}

	if GraphTimeBucketGranularityDay != "day" {
		t.Fatalf("unexpected time bucket granularity: %s", GraphTimeBucketGranularityDay)
	}
}
