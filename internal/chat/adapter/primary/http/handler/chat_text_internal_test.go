package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeUIActionQuickAdd(t *testing.T) {
	t.Run("normalizes trim and draft_cancel operation fallback", func(t *testing.T) {
		requestContext := map[string]interface{}{
			ContextKeyUIAction: map[string]interface{}{
				ContextKeyUIActionType: "draft_cancel",
				ContextKeyQuickAdd: map[string]interface{}{
					"contract_version": " quick-add-v1 ",
					"entity":           "tag",
					"idempotency_key":  "  qa-2  ",
				},
			},
		}

		normalized := normalizeUIActionQuickAdd(requestContext)
		uiAction := normalized[ContextKeyUIAction].(map[string]interface{})
		quickAdd := uiAction[ContextKeyQuickAdd].(map[string]interface{})

		require.Equal(t, "quick-add-v1", quickAdd["contract_version"])
		require.Equal(t, "tag", quickAdd["entity"])
		require.Equal(t, "cancel", quickAdd["operation"])
		require.Equal(t, "qa-2", quickAdd["idempotency_key"])
	})

	t.Run("drops malformed quick_add object", func(t *testing.T) {
		requestContext := map[string]interface{}{
			ContextKeyUIAction: map[string]interface{}{
				ContextKeyUIActionType: "draft_accept",
				ContextKeyQuickAdd:     "invalid",
			},
		}

		normalized := normalizeUIActionQuickAdd(requestContext)
		uiAction := normalized[ContextKeyUIAction].(map[string]interface{})
		_, exists := uiAction[ContextKeyQuickAdd]
		require.False(t, exists)
	})
}
