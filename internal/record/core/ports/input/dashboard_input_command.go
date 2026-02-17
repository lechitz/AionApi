package input

import "time"

// UpsertMetricDefinitionCommand contains input data to create/update a metric definition.
type UpsertMetricDefinitionCommand struct {
	ID          *uint64
	MetricKey   string
	DisplayName string
	CategoryID  *uint64
	TagID       uint64
	TagIDs      []uint64
	ValueSource string
	Aggregation string
	Unit        string
	GoalDefault *float64
	IsActive    *bool
}

// UpsertGoalTemplateCommand contains input data to create/update a goal template.
type UpsertGoalTemplateCommand struct {
	ID          *uint64
	MetricKey   string
	Title       string
	TargetValue float64
	Comparison  string
	Period      string
	IsActive    *bool
}

// DashboardSnapshotQuery contains input parameters for dashboard snapshot queries.
type DashboardSnapshotQuery struct {
	Date     time.Time
	Timezone string
}
