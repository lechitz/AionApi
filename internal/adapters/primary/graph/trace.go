package graph

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"go.opentelemetry.io/otel/attribute"
)

// TraceAttributesFromCategory generates OTEL trace attributes from a category DTO.
func TraceAttributesFromCategory(category model.DtoCreateCategory) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("category_name", category.Name),
	}
	if category.Description != nil {
		attrs = append(attrs, attribute.String("category_description", *category.Description))
	}
	if category.ColorHex != nil {
		attrs = append(attrs, attribute.String("category_color", *category.ColorHex))
	}
	if category.Icon != nil {
		attrs = append(attrs, attribute.String("category_icon", *category.Icon))
	}
	return attrs
}
