package usecase

import (
	"context"
	"strings"
	"time"

	auditdomain "github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/chat/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

func (s *ChatService) persistAuditActionEvent(
	ctx context.Context,
	userID uint64,
	requestContext map[string]interface{},
	result *domain.ChatResult,
) {
	if s.auditService == nil || requestContext == nil || result == nil {
		return
	}

	rawAction, ok := requestContext[ContextKeyUIAction].(map[string]interface{})
	if !ok || rawAction == nil {
		return
	}

	uiActionType := readString(rawAction, ContextKeyUIActionType)
	draftID := readString(rawAction, ContextKeyDraftID)
	if uiActionType == "" || draftID == "" {
		return
	}

	consentRequired, consentConfirmed, consentPolicyVersion := extractConsent(rawAction)
	quickAddContractVersion, quickAddEntity, quickAddOperation, quickAddIdempotencyKey := extractQuickAdd(rawAction)
	status, entityID, messageCode := extractActionResultFields(result.UI)

	operation := quickAddOperation
	if operation == "" {
		operation = inferOperation(uiActionType)
	}

	event := auditdomain.AuditActionEvent{
		TimestampUTC:            time.Now().UTC(),
		UserID:                  userID,
		Source:                  "aion-api",
		TraceID:                 extractTraceID(ctx),
		RequestID:               extractRequestID(ctx),
		UIActionType:            uiActionType,
		DraftID:                 draftID,
		Action:                  uiActionType,
		Entity:                  defaultString(quickAddEntity, "unknown"),
		Operation:               operation,
		Status:                  status,
		EntityID:                entityID,
		ConsentRequired:         consentRequired,
		ConsentConfirmed:        consentConfirmed,
		ConsentPolicyVersion:    consentPolicyVersion,
		QuickAddContractVersion: quickAddContractVersion,
		QuickAddIdempotencyKey:  quickAddIdempotencyKey,
		MessageCode:             messageCode,
		PayloadRedacted: map[string]interface{}{
			ContextKeyUIActionType: uiActionType,
			ContextKeyDraftID:      draftID,
			"status":               status,
			"consent_required":     consentRequired,
			"consent_confirmed":    consentConfirmed,
		},
	}

	if err := s.auditService.WriteEvent(ctx, event); err != nil {
		s.logger.ErrorwCtx(ctx, LogAuditEventSaveFailed,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
			LogKeyUIActionType, uiActionType,
			LogKeyDraftID, draftID,
		)
		return
	}

	s.logger.InfowCtx(ctx, LogAuditEventSaved,
		LogKeyUserID, userID,
		LogKeyUIActionType, uiActionType,
		LogKeyDraftID, draftID,
		LogKeyTraceID, event.TraceID,
	)
}

func extractConsent(rawAction map[string]interface{}) (bool, bool, string) {
	rawConsent, ok := rawAction[ContextKeyConsent].(map[string]interface{})
	if !ok || rawConsent == nil {
		return false, false, ""
	}
	required, _ := rawConsent["required"].(bool)
	confirmed, _ := rawConsent["confirmed"].(bool)
	policyVersion := readString(rawConsent, "policy_version")
	return required, confirmed, policyVersion
}

func extractQuickAdd(rawAction map[string]interface{}) (string, string, string, string) {
	rawQuickAdd, ok := rawAction[ContextKeyQuickAdd].(map[string]interface{})
	if !ok || rawQuickAdd == nil {
		return "", "", "", ""
	}
	return readString(rawQuickAdd, "contract_version"),
		readString(rawQuickAdd, "entity"),
		readString(rawQuickAdd, "operation"),
		readString(rawQuickAdd, "idempotency_key")
}

func extractActionResultFields(ui map[string]interface{}) (string, string, string) {
	defaultStatus := "success"
	if ui == nil {
		return defaultStatus, "", ""
	}
	actionResult, ok := ui[ContextKeyActionResult].(map[string]interface{})
	if !ok || actionResult == nil {
		return defaultStatus, "", ""
	}
	status := readString(actionResult, "status")
	if status == "" {
		status = defaultStatus
	}
	entityID := readString(actionResult, "entity_id")
	messageCode := readString(actionResult, "message_code")
	return status, entityID, messageCode
}

func extractTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxkeys.TraceID).(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func extractRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxkeys.RequestID).(string); ok {
		return strings.TrimSpace(v)
	}
	return ""
}

func inferOperation(uiActionType string) string {
	switch strings.TrimSpace(uiActionType) {
	case "draft_cancel":
		return "cancel"
	default:
		return "create"
	}
}

func readString(data map[string]interface{}, key string) string {
	value, _ := data[key].(string)
	return strings.TrimSpace(value)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
