package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/AionApi/internal/audit/core/domain"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/roles"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type listEventsResponse struct {
	Items []eventResponse `json:"items"`
	Count int             `json:"count"`
}

type eventResponse struct {
	EventID      string    `json:"event_id"`
	TimestampUTC time.Time `json:"timestamp_utc"`
	UserID       uint64    `json:"user_id"`
	TraceID      string    `json:"trace_id"`
	RequestID    string    `json:"request_id,omitempty"`
	UIActionType string    `json:"ui_action_type"`
	DraftID      string    `json:"draft_id"`
	Action       string    `json:"action"`
	Entity       string    `json:"entity"`
	Operation    string    `json:"operation"`
	Status       string    `json:"status"`
	EntityID     string    `json:"entity_id,omitempty"`
	MessageCode  string    `json:"message_code,omitempty"`
}

// ListEvents handles GET /audit/events for authenticated users.
func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuditHandler).Start(r.Context(), SpanListEventsHandler)
	defer span.End()

	userIDValue := ctx.Value(ctxkeys.UserID)
	if userIDValue == nil {
		httpresponse.WriteAuthErrorSpan(ctx, w, span, sharederrors.NewAuthenticationError(errMissingUserID), h.Logger)
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		httpresponse.WriteAuthErrorSpan(ctx, w, span, sharederrors.NewAuthenticationError(errInvalidUserID), h.Logger)
		return
	}

	filter, err := buildFilter(ctx, r, userID)
	if err != nil {
		if strings.Contains(err.Error(), errForbiddenCrossUser) {
			httpresponse.WriteDomainErrorSpan(ctx, w, span, sharederrors.ErrForbidden(errForbiddenCrossUser), errForbiddenCrossUser, h.Logger)
			return
		}
		httpresponse.WriteValidationErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	events, err := h.Service.ListEvents(ctx, filter)
	if err != nil {
		logCrossUserAuditQuery(ctx, h, userID, filter, 0, err)
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, errAuditService, h.Logger)
		return
	}

	logCrossUserAuditQuery(ctx, h, userID, filter, len(events), nil)

	response := listEventsResponse{Items: toResponseItems(events), Count: len(events)}
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int("audit.count", len(events)),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
	)
	span.SetStatus(codes.Ok, msgAuditEventsListed)

	httpresponse.WriteSuccess(w, http.StatusOK, response, msgAuditEventsListed)
}

func buildFilter(ctx context.Context, r *http.Request, userID uint64) (domain.AuditActionEventFilter, error) {
	query := r.URL.Query()
	filterUserID, err := resolveFilterUserID(ctx, query.Get(queryUserID), userID)
	if err != nil {
		return domain.AuditActionEventFilter{}, err
	}

	filter := domain.AuditActionEventFilter{
		UserID:   &filterUserID,
		TraceID:  strings.TrimSpace(query.Get(queryTraceID)),
		DraftID:  strings.TrimSpace(query.Get(queryDraftID)),
		Statuses: query[queryStatus],
		Limit:    defaultLimit,
		Offset:   0,
	}

	if err := applyPagination(query.Get(queryLimit), query.Get(queryOffset), &filter); err != nil {
		return domain.AuditActionEventFilter{}, err
	}
	if err := applyTimeRange(query.Get(queryFromUTC), query.Get(queryToUTC), &filter); err != nil {
		return domain.AuditActionEventFilter{}, err
	}

	if filter.FromUTC != nil && filter.ToUTC != nil && filter.FromUTC.After(*filter.ToUTC) {
		return domain.AuditActionEventFilter{}, sharederrors.NewValidationError(queryFromUTC, "must be before or equal to to_utc")
	}

	return filter, nil
}

func resolveFilterUserID(ctx context.Context, rawUserID string, defaultUserID uint64) (uint64, error) {
	trimmed := strings.TrimSpace(rawUserID)
	if trimmed == "" {
		return defaultUserID, nil
	}
	if !hasAdminRole(extractRolesFromContext(ctx)) {
		return 0, sharederrors.ErrForbidden(errForbiddenCrossUser)
	}

	parsedUserID, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil || parsedUserID == 0 {
		return 0, sharederrors.NewValidationError(queryUserID, "must be a positive integer")
	}
	return parsedUserID, nil
}

func applyPagination(rawLimit, rawOffset string, filter *domain.AuditActionEventFilter) error {
	trimmedLimit := strings.TrimSpace(rawLimit)
	if trimmedLimit != "" {
		limit, err := strconv.Atoi(trimmedLimit)
		if err != nil || limit <= 0 {
			return sharederrors.NewValidationError(queryLimit, "must be a positive integer")
		}
		if limit > maxLimit {
			limit = maxLimit
		}
		filter.Limit = limit
	}

	trimmedOffset := strings.TrimSpace(rawOffset)
	if trimmedOffset == "" {
		return nil
	}

	offset, err := strconv.Atoi(trimmedOffset)
	if err != nil || offset < 0 {
		return sharederrors.NewValidationError(queryOffset, "must be a non-negative integer")
	}
	filter.Offset = offset
	return nil
}

func applyTimeRange(rawFrom, rawTo string, filter *domain.AuditActionEventFilter) error {
	trimmedFrom := strings.TrimSpace(rawFrom)
	if trimmedFrom != "" {
		from, err := time.Parse(time.RFC3339, trimmedFrom)
		if err != nil {
			return sharederrors.NewValidationError(queryFromUTC, "must be RFC3339")
		}
		filter.FromUTC = &from
	}

	trimmedTo := strings.TrimSpace(rawTo)
	if trimmedTo == "" {
		return nil
	}

	to, err := time.Parse(time.RFC3339, trimmedTo)
	if err != nil {
		return sharederrors.NewValidationError(queryToUTC, "must be RFC3339")
	}
	filter.ToUTC = &to
	return nil
}

func extractRolesFromContext(ctx context.Context) []string {
	claimsVal := ctx.Value(ctxkeys.Claims)
	if claimsVal == nil {
		return []string{}
	}

	claims, ok := claimsVal.(map[string]any)
	if !ok {
		return []string{}
	}

	rolesVal, exists := claims[commonkeys.Roles]
	if !exists {
		return []string{}
	}
	if typedRoles, ok := rolesVal.([]string); ok {
		return typedRoles
	}
	if rolesAny, ok := rolesVal.([]any); ok {
		typedRoles := make([]string, 0, len(rolesAny))
		for _, value := range rolesAny {
			roleName, ok := value.(string)
			if ok {
				typedRoles = append(typedRoles, roleName)
			}
		}
		return typedRoles
	}
	return []string{}
}

func hasAdminRole(rolesList []string) bool {
	for _, roleName := range rolesList {
		if roleName == roles.Admin {
			return true
		}
	}
	return false
}

func logCrossUserAuditQuery(
	ctx context.Context,
	h *Handler,
	actorUserID uint64,
	filter domain.AuditActionEventFilter,
	resultCount int,
	listErr error,
) {
	if filter.UserID == nil || *filter.UserID == actorUserID {
		return
	}

	targetUserID := *filter.UserID
	if listErr != nil {
		h.Logger.ErrorwCtx(ctx, MsgCrossUserAuditQueryFailed,
			commonkeys.Error, listErr.Error(),
			"actor_user_id", actorUserID,
			"target_user_id", targetUserID,
			queryTraceID, filter.TraceID,
			queryDraftID, filter.DraftID,
			"limit", filter.Limit,
			"offset", filter.Offset,
		)
		return
	}

	h.Logger.InfowCtx(ctx, MsgCrossUserAuditQuery,
		"actor_user_id", actorUserID,
		"target_user_id", targetUserID,
		queryTraceID, filter.TraceID,
		queryDraftID, filter.DraftID,
		"status_count", len(filter.Statuses),
		"limit", filter.Limit,
		"offset", filter.Offset,
		"result_count", resultCount,
	)
}

func toResponseItems(events []domain.AuditActionEvent) []eventResponse {
	items := make([]eventResponse, 0, len(events))
	for _, event := range events {
		items = append(items, eventResponse{
			EventID:      event.EventID,
			TimestampUTC: event.TimestampUTC,
			UserID:       event.UserID,
			TraceID:      event.TraceID,
			RequestID:    event.RequestID,
			UIActionType: event.UIActionType,
			DraftID:      event.DraftID,
			Action:       event.Action,
			Entity:       event.Entity,
			Operation:    event.Operation,
			Status:       event.Status,
			EntityID:     event.EntityID,
			MessageCode:  event.MessageCode,
		})
	}
	return items
}
