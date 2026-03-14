package usecase

import (
	"fmt"
	"sort"
	"strconv"
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

	nodes, edges, addNode, addEdge := newGraphProjectionAccumulator()

	userNodeID := graphNodeID(recorddomain.GraphNodePrefixUser, strconv.FormatUint(in.UserID, 10))
	addUserNode(in.UserID, userNodeID, addNode)

	categoryIDs := addCategoryNodes(in, addNode)
	tagIDs := addTagNodesAndEdges(in, categoryIDs, addNode, addEdge)
	recordIDs := addRecordNodesAndEdges(in, loc, userNodeID, tagIDs, addNode, addEdge)
	addInsightNodesAndEdges(in, recordIDs, tagIDs, categoryIDs, addNode, addEdge)

	return recorddomain.GraphProjection{
		Version:     recorddomain.GraphProjectionVersionV1,
		UserID:      in.UserID,
		GeneratedAt: generatedAt.UTC(),
		Nodes:       sortedGraphNodes(nodes),
		Edges:       sortedGraphEdges(edges),
	}
}

func newGraphProjectionAccumulator() (
	map[string]recorddomain.GraphNode,
	map[string]recorddomain.GraphEdge,
	func(recorddomain.GraphNode),
	func(recorddomain.GraphEdge),
) {
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

	return nodes, edges, addNode, addEdge
}

func addUserNode(
	userID uint64,
	userNodeID string,
	addNode func(recorddomain.GraphNode),
) {
	addNode(recorddomain.GraphNode{
		ID:     userNodeID,
		Type:   recorddomain.GraphNodeTypeUser,
		Label:  fmt.Sprintf("User %d", userID),
		UserID: userID,
		Metadata: map[string]string{
			"user_id": strconv.FormatUint(userID, 10),
		},
	})
}

func addCategoryNodes(
	in GraphProjectionBuildInput,
	addNode func(recorddomain.GraphNode),
) map[uint64]string {
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
				"category_id": strconv.FormatUint(category.ID, 10),
				"name":        category.Name,
				"description": category.Description,
				"color":       category.Color,
				"icon":        category.Icon,
			}),
		})
	}
	return categoryIDs
}

func addTagNodesAndEdges(
	in GraphProjectionBuildInput,
	categoryIDs map[uint64]string,
	addNode func(recorddomain.GraphNode),
	addEdge func(recorddomain.GraphEdge),
) map[uint64]string {
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
				"tag_id":      strconv.FormatUint(tag.ID, 10),
				"name":        tag.Name,
				"description": tag.Description,
				"icon":        tag.Icon,
				"category_id": strconv.FormatUint(tag.CategoryID, 10),
				"usage_count": strconv.Itoa(tag.UsageCount),
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
					"tag_id":      strconv.FormatUint(tag.ID, 10),
					"category_id": strconv.FormatUint(tag.CategoryID, 10),
				},
			})
		}
	}
	return tagIDs
}

func addRecordNodesAndEdges(
	in GraphProjectionBuildInput,
	loc *time.Location,
	userNodeID string,
	tagIDs map[uint64]string,
	addNode func(recorddomain.GraphNode),
	addEdge func(recorddomain.GraphEdge),
) map[uint64]string {
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
				"record_id":        strconv.FormatUint(record.ID, 10),
				"tag_id":           strconv.FormatUint(record.TagID, 10),
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
				"user_id":   strconv.FormatUint(in.UserID, 10),
				"record_id": strconv.FormatUint(record.ID, 10),
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
					"record_id": strconv.FormatUint(record.ID, 10),
					"tag_id":    strconv.FormatUint(record.TagID, 10),
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
				"record_id": strconv.FormatUint(record.ID, 10),
				"date":      eventLocal.Format("2006-01-02"),
				"timezone":  loc.String(),
			},
		})
	}
	return recordIDs
}

func addInsightNodesAndEdges(
	in GraphProjectionBuildInput,
	recordIDs map[uint64]string,
	tagIDs map[uint64]string,
	categoryIDs map[uint64]string,
	addNode func(recorddomain.GraphNode),
	addEdge func(recorddomain.GraphEdge),
) {
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
				"confidence":         strconv.Itoa(insight.Confidence),
				"summary":            insight.Summary,
				"recommended_action": derefString(insight.RecommendedAction),
				"generated_at_utc":   insight.GeneratedAt.UTC().Format(time.RFC3339),
				"metric_keys":        strings.Join(insight.MetricKeys, ","),
			}),
		})

		addInsightSupportedRecordEdges(in, insight.ID, nodeID, recordIDs, addEdge)
		addInsightScopedTagEdges(in, insight.ID, nodeID, tagIDs, addEdge)
		addInsightScopedCategoryEdges(in, insight.ID, nodeID, categoryIDs, addEdge)
	}
}

func addInsightSupportedRecordEdges(
	in GraphProjectionBuildInput,
	insightID string,
	nodeID string,
	recordIDs map[uint64]string,
	addEdge func(recorddomain.GraphEdge),
) {
	for _, recordID := range in.InsightSupportedRecordIDs[insightID] {
		if recordNodeID, ok := recordIDs[recordID]; ok {
			addEdge(recorddomain.GraphEdge{
				ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightSupportedByRecord), nodeID, recordNodeID),
				Type:   recorddomain.GraphEdgeTypeInsightSupportedByRecord,
				From:   nodeID,
				To:     recordNodeID,
				UserID: in.UserID,
				Metadata: map[string]string{
					"insight_id": insightID,
					"record_id":  strconv.FormatUint(recordID, 10),
				},
			})
		}
	}
}

func addInsightScopedTagEdges(
	in GraphProjectionBuildInput,
	insightID string,
	nodeID string,
	tagIDs map[uint64]string,
	addEdge func(recorddomain.GraphEdge),
) {
	for _, tagID := range in.InsightScopedTagIDs[insightID] {
		if tagNodeID, ok := tagIDs[tagID]; ok {
			addEdge(recorddomain.GraphEdge{
				ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightScopedToTag), nodeID, tagNodeID),
				Type:   recorddomain.GraphEdgeTypeInsightScopedToTag,
				From:   nodeID,
				To:     tagNodeID,
				UserID: in.UserID,
				Metadata: map[string]string{
					"insight_id": insightID,
					"tag_id":     strconv.FormatUint(tagID, 10),
				},
			})
		}
	}
}

func addInsightScopedCategoryEdges(
	in GraphProjectionBuildInput,
	insightID string,
	nodeID string,
	categoryIDs map[uint64]string,
	addEdge func(recorddomain.GraphEdge),
) {
	for _, categoryID := range in.InsightScopedCategoryIDs[insightID] {
		if categoryNodeID, ok := categoryIDs[categoryID]; ok {
			addEdge(recorddomain.GraphEdge{
				ID:     graphEdgeID(string(recorddomain.GraphEdgeTypeInsightScopedToCategory), nodeID, categoryNodeID),
				Type:   recorddomain.GraphEdgeTypeInsightScopedToCategory,
				From:   nodeID,
				To:     categoryNodeID,
				UserID: in.UserID,
				Metadata: map[string]string{
					"insight_id":  insightID,
					"category_id": strconv.FormatUint(categoryID, 10),
				},
			})
		}
	}
}

func sortedGraphNodes(nodes map[string]recorddomain.GraphNode) []recorddomain.GraphNode {
	outNodes := make([]recorddomain.GraphNode, 0, len(nodes))
	for _, node := range nodes {
		outNodes = append(outNodes, node)
	}
	sort.Slice(outNodes, func(i, j int) bool {
		return outNodes[i].ID < outNodes[j].ID
	})
	return outNodes
}

func sortedGraphEdges(edges map[string]recorddomain.GraphEdge) []recorddomain.GraphEdge {
	outEdges := make([]recorddomain.GraphEdge, 0, len(edges))
	for _, edge := range edges {
		outEdges = append(outEdges, edge)
	}
	sort.Slice(outEdges, func(i, j int) bool {
		return outEdges[i].ID < outEdges[j].ID
	})
	return outEdges
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
	return graphNodeID(recorddomain.GraphNodePrefixCategory, strconv.FormatUint(categoryID, 10))
}

func tagNodeID(tagID uint64) string {
	return graphNodeID(recorddomain.GraphNodePrefixTag, strconv.FormatUint(tagID, 10))
}

func recordNodeID(recordID uint64) string {
	return graphNodeID(recorddomain.GraphNodePrefixRecord, strconv.FormatUint(recordID, 10))
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
