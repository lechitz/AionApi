package domain

import "time"

// InsightWindow identifies the supported analysis windows for v1 insights/series.
type InsightWindow string

const (
	// InsightWindow7D is the default 7-day insight and analytics window.
	InsightWindow7D InsightWindow = "WINDOW_7D"
	// InsightWindow30D is the 30-day insight and analytics window.
	InsightWindow30D InsightWindow = "WINDOW_30D"
	// InsightWindow90D is the 90-day insight and analytics window.
	InsightWindow90D InsightWindow = "WINDOW_90D"
)

// InsightEvidence explains one supporting fact behind an insight.
type InsightEvidence struct {
	Label string
	Value string
	Kind  string
}

// InsightCard is the canonical explainable insight payload for v1.
type InsightCard struct {
	ID                string
	Type              string
	Title             string
	Summary           string
	Status            string
	Window            InsightWindow
	Confidence        int
	MetricKeys        []string
	RecommendedAction *string
	Evidence          []InsightEvidence
	GeneratedAt       time.Time
}

// AnalyticsPoint represents one point in an analytics series.
type AnalyticsPoint struct {
	Timestamp time.Time
	Value     *float64
	Label     *string
}

// AnalyticsSeriesResult is a compact time-series payload for dashboard consumers.
type AnalyticsSeriesResult struct {
	SeriesKey string
	Window    InsightWindow
	Points    []AnalyticsPoint
	Summary   *string
}

// MetricDefinition configures how a dashboard metric is computed from records.
type MetricDefinition struct {
	ID          uint64
	UserID      uint64
	MetricKey   string
	DisplayName string
	CategoryID  *uint64
	TagID       uint64
	TagIDs      []uint64
	ValueSource string
	Aggregation string
	Unit        string
	GoalDefault *float64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GoalTemplate configures a deterministic daily goal bound to a metric key.
type GoalTemplate struct {
	ID          uint64
	UserID      uint64
	MetricKey   string
	Title       string
	TargetValue float64
	Comparison  string
	Period      string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const (
	// DashboardChecklistModeCountGoal represents a checklist backed by a discrete count with an explicit target.
	DashboardChecklistModeCountGoal = "count_goal"
	// DashboardChecklistModeCountOnly represents a checklist backed by a discrete count without a configured target.
	DashboardChecklistModeCountOnly = "count_only"
	// DashboardChecklistModeBinaryDone represents a checklist backed by a latest-state boolean completion model.
	DashboardChecklistModeBinaryDone = "binary_done"
)

// DashboardChecklistValue is the explicit checklist-oriented semantic payload for dashboard consumers.
type DashboardChecklistValue struct {
	MetricKey        string
	Label            string
	CompletedCount   int
	TargetCount      *int
	RemainingCount   *int
	CompletionRatio  float64
	Status           string
	Mode             string
}

// DashboardMetricValue represents a computed metric value for a given date.
type DashboardMetricValue struct {
	MetricKey   string
	Label       string
	Value       float64
	Unit        string
	Target      *float64
	ProgressPct float64
	Status      string
	Checklist   *DashboardChecklistValue
}

// DashboardGoalValue represents daily goal progress.
type DashboardGoalValue struct {
	GoalID      uint64
	Title       string
	MetricKey   string
	Current     float64
	Target      float64
	ProgressPct float64
	Status      string
}

// DashboardSnapshot is the aggregate payload consumed by /dashboard.
type DashboardSnapshot struct {
	Date     time.Time
	Timezone string
	Metrics  []DashboardMetricValue
	Goals    []DashboardGoalValue
}

// Dashboard widget supported sizes.
const (
	DashboardWidgetSizeSmall  = "small"
	DashboardWidgetSizeMedium = "medium"
	DashboardWidgetSizeLarge  = "large"

	DashboardWidgetTypeKPINumber    = "kpi_number"
	DashboardWidgetTypeGoalProgress = "goal_progress"
	DashboardWidgetTypeTrendLine    = "trend_line"
	DashboardWidgetTypeChecklist    = "checklist"
	DashboardWidgetTypeInsightFeed  = "insight_feed"
)

// MaxLargeWidgetsPerDashboard limits how many large widgets can exist per view.
const (
	MaxLargeWidgetsPerDashboard = 3
)

// DashboardWidget configures one dashboard card/widget instance in a view.
type DashboardWidget struct {
	ID                 uint64
	UserID             uint64
	ViewID             uint64
	MetricDefinitionID uint64
	WidgetType         string
	Size               string
	OrderIndex         int
	TitleOverride      *string
	ConfigJSON         string
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// DashboardView is a user-defined dashboard layout that contains widgets.
type DashboardView struct {
	ID        uint64
	UserID    uint64
	Name      string
	IsDefault bool
	Widgets   []DashboardWidget
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MetricDefinitionSuggestion is a deterministic proposal users can accept before creation.
type MetricDefinitionSuggestion struct {
	MetricKey   string
	DisplayName string
	CategoryID  *uint64
	TagIDs      []uint64
	ValueSource string
	Aggregation string
	Unit        string
	Reason      string
}
