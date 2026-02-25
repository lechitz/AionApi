package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractUIActionMetadata(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		actionType, draftID := extractUIActionMetadata(nil)
		require.Equal(t, "", actionType)
		require.Equal(t, "", draftID)
	})

	t.Run("missing ui_action", func(t *testing.T) {
		actionType, draftID := extractUIActionMetadata(map[string]interface{}{"k": "v"})
		require.Equal(t, "", actionType)
		require.Equal(t, "", draftID)
	})

	t.Run("ui_action malformed type", func(t *testing.T) {
		actionType, draftID := extractUIActionMetadata(map[string]interface{}{
			ContextKeyUIAction: "invalid",
		})
		require.Equal(t, "", actionType)
		require.Equal(t, "", draftID)
	})

	t.Run("valid metadata", func(t *testing.T) {
		actionType, draftID := extractUIActionMetadata(map[string]interface{}{
			ContextKeyUIAction: map[string]interface{}{
				ContextKeyUIActionType: "draft_cancel",
				ContextKeyDraftID:      "draft-123",
			},
		})
		require.Equal(t, "draft_cancel", actionType)
		require.Equal(t, "draft-123", draftID)
	})
}
