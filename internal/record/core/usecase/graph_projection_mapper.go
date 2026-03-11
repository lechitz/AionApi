package usecase

import (
	"fmt"
	"sort"
	"strings"
	"time"

	categorydomain "github.com/lechitz/AionApi/internal/category/core/domain"
	recorddomain "github.com/lechitz/AionApi/internal/record/core/domain"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
)

// GraphProjectionBuildInput contains the canonical entities used to derive the graph projection.
type GraphProjectionBuildInput struct {
	UserID                    uint64
	GeneratedAt               time.Time
	Timezone                  string
	Records                   []recorddomain.Record
	Tags                      []tagdomain.Tag
	Categories                []categorydomain.Category
	Insights                  []recorddomain.InsightCard
	InsightSupportedRecordIDs map[string][]uint64
	InsightScopedTagIDs       map[string][]uint64
	InsightScopedCategoryIDs  map[string][]uint64
}

// BuildGraphProjection derives a stable graph-projection-v1 payload from canonical Aion entities.
func BuildGraphProjection(in GraphProjectionBuildInput) recorddomain.GraphProjection {
	generatedAt := in.GeneratedAt
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}

	loc, _ := resolveInsightLocation(in.Timezone)
	if loc == nil {
		loc = time.UTC
	}

	nodes := make(map[string]recorddomain.GraphNode)
	edges := make(map[string]recorddomain.GraphEdge)

	addNode := func(node recorddomain.GraphNode) {
		if node.ID == "" {
			return
		}
		if _, exists := nodes[node.ID]; exists {
			return
		}
		nodes[node.ID] = node
	}

	addEdge := func(edge recorddomain.GraphEdge) {
		if edge.ID == "" {
			return
		}
		if _, exists := edges[edge.ID]; exists {
			return
		}
		edges[edge.ID] = edge
	}

	userNodeID := graphNodeID(recorddomain.GraphNodePrefixUser, fmt.Sprintf("%d", in.UserID))
	addNode(recorddomain.GraphNode{
		ID:     userNodeID,
		Type:   recorddomain.GraphNodeTypeUser,
		Label:  fmt.Sprintf("User %d", in.UserID),
		UserID: in.UserID,
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", in.UserID),
		},
	})

	categoryIDs := make(map[uint64]string, len(in.Categories))
	for _, category := range in.Categories {
		nodeID := categoryNodeID(category.ID)
		categoryIDs[category.ID] = nodeID
		addNode(recorddomain.GraphNode{
			ID:     nodeID,
			Type:   recorddomain.GraphNodeTypeCategory,
			Label:  category.Name,
			UserID: in.UserID,
			Metadata: compactMetadata(map[string]string{
				"category_id": fmt.Sprintf("%d", category.ID),
				"name":        category.Name,
				"description": category.Description,
				"color":       category.Color,
				"icon":        category.Icon,
			}),
		})
	}

	tagIDs := make(map[uint64]string, len(in.Tags))
	for _, tag := range in.Tags {
		nodeID := tagNodeID(tag.ID)
		tagIDs[tag.ID] = nodeID
		addNode(recorddomain.GraphNode{
			ID:     nodeID,
			Type:   recorddomain.GraphNodeTypeTag,
			Label:  tag.Name,
			UserID: in.UserID,
			Metadata: compactMetadata(map[string]string{
				"tag_id":      fmt.Sprintf("%d", tag.ID),
				"name":        tag.Name,
				"description": tag.Description,
				"icon":        tag.Icon,
				"category_id": fmt.Sprintf("%d", tag.CategoryID),
				"usage_count": fmt.Sprintf("%d", tag.UsageCount),
			}),
		})
		if categoryNodeID, ok := categoryIDs[tag.CategoryID]; ok {
			addEdge(recorddomain.GraphEdge{
				ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeTagBelongsToCategory), nodeID, categoryNodeID),
				Type:   recorddomain.GraphEdgeTypeTagBelongsToCategory,
				From:   nodeID,
				To:     categoryNodeID,
				UserID: in.UserID,
				Metadata: map[string]string{
					"tag_id":      fmt.Sprintf("%d", tag.ID),
					"category_id": fmt.Sprintf("%d", tag.CategoryID),
				},
			})
		}
	}

	recordIDs := make(map[uint64]string, len(in.Records))
	for _, record := range in.Records {
		nodeID := recordNodeID(record.ID)
		recordIDs[record.ID] = nodeID
		eventLocal := record.EventTime.In(loc)
		timeBucketID := timeBucketNodeID(eventLocal)

		addNode(recorddomain.GraphNode{
			ID:     nodeID,
			Type:   recorddomain.GraphNodeTypeRecord,
			Label:  fmt.Sprintf("Record %d", record.ID),
			UserID: in.UserID,
			Metadata: compactMetadata(map[string]string{
				"record_id":        fmt.Sprintf("%d", record.ID),
				"tag_id":           fmt.Sprintf("%d", record.TagID),
				"event_time_utc":   record.EventTime.UTC().Format(time.RFC3339),
				"event_date_local": eventLocal.Format("2006-01-02"),
				"timezone":         loc.String(),
				"status":           derefString(record.Status),
				"source":           derefString(record.Source),
				"description":      derefString(record.Description),
			}),
		})

		addNode(recorddomain.GraphNode{
			ID:     timeBucketID,
			Type:   recorddomain.GraphNodeTypeTimeBucket,
			Label:  eventLocal.Format("2006-01-02"),
			UserID: in.UserID,
			Metadata: compactMetadata(map[string]string{
				"granularity": string(recorddomain.GraphTimeBucketGranularityDay),
				"date":        eventLocal.Format("2006-01-02"),
				"timezone":    loc.String(),
			}),
		})

		addEdge(recorddomain.GraphEdge{
			ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeUserCreatedRecord), userNodeID, nodeID),
			Type:   recorddomain.GraphEdgeTypeUserCreatedRecord,
			From:   userNodeID,
			To:     nodeID,
			UserID: in.UserID,
			Metadata: map[string]string{
				"user_id":   fmt.Sprintf("%d", in.UserID),
				"record_id": fmt.Sprintf("%d", record.ID),
			},
		})

		if tagNodeID, ok := tagIDs[record.TagID]; ok {
			addEdge(recorddomain.GraphEdge{
				ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeRecordTaggedAs), nodeID, tagNodeID),
				Type:   recorddomain.GraphEdgeTypeRecordTaggedAs,
				From:   nodeID,
				To:     tagNodeID,
				UserID: in.UserID,
				Metadata: map[string]string{
					"record_id": fmt.Sprintf("%d", record.ID),
					"tag_id":    fmt.Sprintf("%d", record.TagID),
				},
			})
		}

		addEdge(recorddomain.GraphEdge{
			ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeRecordOccurredIn), nodeID, timeBucketID),
			Type:   recorddomain.GraphEdgeTypeRecordOccurredIn,
			From:   nodeID,
			To:     timeBucketID,
			UserID: in.UserID,
			Metadata: map[string]string{
				"record_id": fmt.Sprintf("%d", record.ID),
				"date":      eventLocal.Format("2006-01-02"),
				"timezone":  loc.String(),
			},
		})
	}

	for _, insight := range in.Insights {
		nodeID := insightNodeID(insight.ID)
		addNode(recorddomain.GraphNode{
			ID:     nodeID,
			Type:   recorddomain.GraphNodeTypeInsight,
			Label:  insight.Title,
			UserID: in.UserID,
			Metadata: compactMetadata(map[string]string{
				"insight_id":         insight.ID,
				"type":               insight.Type,
				"status":             insight.Status,
				"window":             string(insight.Window),
				"confidence":         fmt.Sprintf("%d", insight.Confidence),
				"summary":            insight.Summary,
				"recommended_action": derefString(insight.RecommendedAction),
				"generated_at_utc":   insight.GeneratedAt.UTC().Format(time.RFC3339),
				"metric_keys":        strings.Join(insight.MetricKeys, ","),
			}),
		})

		for _, recordID := range in.InsightSupportedRecordIDs[insight.ID] {
			if recordNodeID, ok := recordIDs[recordID]; ok {
				addEdge(recorddomain.GraphEdge{
					ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightSupportedByRecord), nodeID, recordNodeID),
					Type:   recorddomain.GraphEdgeTypeInsightSupportedByRecord,
					From:   nodeID,
					To:     recordNodeID,
					UserID: in.UserID,
					Metadata: map[string]string{
						"insight_id": insight.ID,
						"record_id":  fmt.Sprintf("%d", recordID),
					},
				})
			}
		}

		for _, tagID := range in.InsightScopedTagIDs[insight.ID] {
			if tagNodeID, ok := tagIDs[tagID]; ok {
				addEdge(recorddomain.GraphEdge{
					ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightScopedToTag), nodeID, tagNodeID),
					Type:   recorddomain.GraphEdgeTypeInsightScopedToTag,
					From:   nodeID,
					To:     tagNodeID,
					UserID: in.UserID,
					Metadata: map[string]string{
						"insight_id": insight.ID,
						"tag_id":     fmt.Sprintf("%d", tagID),
					},
				})
			}
		}

		for _, categoryID := range in.InsightScopedCategoryIDs[insight.ID] {
			if categoryNodeID, ok := categoryIDs[categoryID]; ok {
				addEdge(recorddomain.GraphEdge{
					ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightScopedToCategory), nodeID, categoryNodeID),
					Type:   recorddomain.GraphEdgeTypeInsightScopedToCategory,
					From:   nodeID,
					To:     categoryNodeID,
					UserID: in.UserID,
					Metadata: map[string]string{
						"insight_id":  insight.ID,
						"category_id": fmt.Sprintf("%d", categoryID),
					},
				})
			}
		}
	}

	outNodes := make([]recorddomain.GraphNode, 0, len(nodes))
	for _, node := range nodes {
		outNodes = append(outNodes, node)
	}
	sort.Slice(outNodes, func(i, j int) bool {
		return outNodes[i].ID < outNodes[j].ID
	})

	outEdges := make([]recorddomain.GraphEdge, 0, len(edges))
	for _, edge := range edges {
		outEdges = append(outEdges, edge)
	}
	sort.Slice(outEdges, func(i, j int) bool {
		return outEdges[i].ID < outEdges[j].ID
	})

	return recorddomain.GraphProjection{
		Version:     recorddomain.GraphProjectionVersionV1,
		UserID:      in.UserID,
		GeneratedAt: generatedAt.UTC(),
		Nodes:       outNodes,
		Edges:       outEdges,
	}
}

func compactMetadata(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		if strings.TrimSpace(value) == "" {
			continue
		}
		out[key] = value
	}
	return out
}

func graphNodeID(prefix, value string) string {
	return prefix + ":" + value
}

func graphEdgeID(edgeType, fromID, toID string) string {
	return edgeType + ":" + fromID + ":" + toID
}

func categoryNodeID(categoryID uint64) string {
	return graphNodeID(recorddomain.GraphNodePrefixCategory, fmt.Sprintf("%d", categoryID))
}

func tagNodeID(tagID uint64) string {
	return graphNodeID(recorddomain.GraphNodePrefixTag, fmt.Sprintf("%d", tagID))
}

func recordNodeID(recordID uint64) string {
	return graphNodeID(recorddomain.GraphNodePrefixRecord, fmt.Sprintf("%d", recordID))
}

func insightNodeID(insightID string) string {
	return graphNodeID(recorddomain.GraphNodePrefixInsight, insightID)
}

func timeBucketNodeID(day time.Time) string {
	return graphNodeID(recorddomain.GraphNodePrefixTimeBucket, fmt.Sprintf("%s:%s", recorddomain.GraphTimeBucketGranularityDay, day.Format("2006-01-02")))
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
