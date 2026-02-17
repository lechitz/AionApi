package model

import "time"

// MetricDefinition maps aion_api.metric_definitions.
type MetricDefinition struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	UserID      uint64    `gorm:"column:user_id;not null"`
	MetricKey   string    `gorm:"column:metric_key;not null"`
	DisplayName string    `gorm:"column:display_name;not null"`
	CategoryID  *uint64   `gorm:"column:category_id"`
	TagID       uint64    `gorm:"column:tag_id;not null"`
	ValueSource string    `gorm:"column:value_source;not null"`
	Aggregation string    `gorm:"column:aggregation;not null"`
	Unit        string    `gorm:"column:unit;not null"`
	GoalDefault *float64  `gorm:"column:goal_default"`
	IsActive    bool      `gorm:"column:is_active;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (MetricDefinition) TableName() string {
	return "aion_api.metric_definitions"
}

// MetricDefinitionTagBinding maps aion_api.metric_definition_tag_bindings.
type MetricDefinitionTagBinding struct {
	ID                 uint64    `gorm:"column:id;primaryKey"`
	UserID             uint64    `gorm:"column:user_id;not null"`
	MetricDefinitionID uint64    `gorm:"column:metric_definition_id;not null"`
	TagID              uint64    `gorm:"column:tag_id;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;not null"`
}

func (MetricDefinitionTagBinding) TableName() string {
	return "aion_api.metric_definition_tag_bindings"
}

// GoalTemplate maps aion_api.goal_templates.
type GoalTemplate struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	UserID      uint64    `gorm:"column:user_id;not null"`
	MetricKey   string    `gorm:"column:metric_key;not null"`
	Title       string    `gorm:"column:title;not null"`
	TargetValue float64   `gorm:"column:target_value;not null"`
	Comparison  string    `gorm:"column:comparison;not null"`
	Period      string    `gorm:"column:period;not null"`
	IsActive    bool      `gorm:"column:is_active;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
}

func (GoalTemplate) TableName() string {
	return "aion_api.goal_templates"
}
