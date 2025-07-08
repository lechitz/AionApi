package graph

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"
	"go.opentelemetry.io/otel/attribute"
)

// TraceAttributesFromCategory generates OTEL trace attributes from a category DTO.
func TraceAttributesFromCategory(category model.DtoCreateCategory) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String(commonkeys.CategoryName, category.Name),
	}
	if category.Description != nil {
		attrs = append(attrs, attribute.String(commonkeys.CategoryDescription, *category.Description))
	}
	if category.ColorHex != nil {
		attrs = append(attrs, attribute.String(commonkeys.CategoryColor, *category.ColorHex))
	}
	if category.Icon != nil {
		attrs = append(attrs, attribute.String(commonkeys.CategoryIcon, *category.Icon))
	}
	return attrs
}
