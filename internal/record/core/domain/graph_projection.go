package domain

import "time"

// GraphProjectionVersion identifies the internal export contract version.
type GraphProjectionVersion string

const (
	// GraphProjectionVersionV1 is the first stable internal graph projection schema.
	GraphProjectionVersionV1 GraphProjectionVersion = "graph-projection-v1"
)

// GraphNodeType identifies the supported node classes in the v1 projection.
type GraphNodeType string

const (
	GraphNodeTypeUser       GraphNodeType = "user"
	GraphNodeTypeRecord     GraphNodeType = "record"
	GraphNodeTypeTag        GraphNodeType = "tag"
	GraphNodeTypeCategory   GraphNodeType = "category"
	GraphNodeTypeTimeBucket GraphNodeType = "time_bucket"
	GraphNodeTypeInsight    GraphNodeType = "insight"
)

// GraphEdgeType identifies the supported relation classes in the v1 projection.
type GraphEdgeType string

const (
	GraphEdgeTypeUserCreatedRecord        GraphEdgeType = "user_created_record"
	GraphEdgeTypeRecordTaggedAs           GraphEdgeType = "record_tagged_as"
	GraphEdgeTypeTagBelongsToCategory     GraphEdgeType = "tag_belongs_to_category"
	GraphEdgeTypeRecordOccurredIn         GraphEdgeType = "record_occurred_in"
	GraphEdgeTypeInsightSupportedByRecord GraphEdgeType = "insight_supported_by_record"
	GraphEdgeTypeInsightScopedToTag       GraphEdgeType = "insight_scoped_to_tag"
	GraphEdgeTypeInsightScopedToCategory  GraphEdgeType = "insight_scoped_to_category"
)

// GraphProjection node id prefixes are stable naming hints for future mappers/exports.
const (
	GraphNodePrefixUser       = "user"
	GraphNodePrefixRecord     = "record"
	GraphNodePrefixTag        = "tag"
	GraphNodePrefixCategory   = "category"
	GraphNodePrefixTimeBucket = "time_bucket"
	GraphNodePrefixInsight    = "insight"
)

// GraphTimeBucketGranularity identifies supported time-bucket shapes for graph exports.
type GraphTimeBucketGranularity string

const (
	GraphTimeBucketGranularityDay GraphTimeBucketGranularity = "day"
)

// GraphNode is one exported node in the graph projection.
//
// Metadata is intentionally string-only in v1 so the export shape remains stable
// and easy to consume from local JSON/CSV/Neo4j lab tooling.
type GraphNode struct {
	ID       string
	Type     GraphNodeType
	Label    string
	UserID   uint64
	Metadata map[string]string
}

// GraphEdge is one exported relation in the graph projection.
type GraphEdge struct {
	ID       string
	Type     GraphEdgeType
	From     string
	To       string
	UserID   uint64
	Metadata map[string]string
}

// GraphProjection is the internal export payload for graph-ready lab usage.
//
// This schema is not a runtime source of truth. It is a derived, export-oriented
// representation built from canonical Aion entities.
type GraphProjection struct {
	Version     GraphProjectionVersion
	UserID      uint64
	GeneratedAt time.Time
	Nodes       []GraphNode
	Edges       []GraphEdge
}
