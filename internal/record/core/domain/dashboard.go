package domain

import "time"

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

// DashboardMetricValue represents a computed metric value for a given date.
type DashboardMetricValue struct {
	MetricKey   string
	Label       string
	Value       float64
	Unit        string
	Target      *float64
	ProgressPct float64
	Status      string
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
