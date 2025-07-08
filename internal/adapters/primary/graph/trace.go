package graph

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/shared/common"
	"go.opentelemetry.io/otel/attribute"
)

// TraceAttributesFromCategory generates OTEL trace attributes from a category DTO.
func TraceAttributesFromCategory(category model.DtoCreateCategory) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String(common.CategoryName, category.Name),
	}
	if category.Description != nil {
		attrs = append(attrs, attribute.String(common.CategoryDescription, *category.Description))
	}
	if category.ColorHex != nil {
		attrs = append(attrs, attribute.String(common.CategoryColor, *category.ColorHex))
	}
	if category.Icon != nil {
		attrs = append(attrs, attribute.String(common.CategoryIcon, *category.Icon))
	}
	return attrs
}
