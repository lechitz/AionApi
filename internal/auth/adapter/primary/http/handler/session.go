package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lechitz/aion-api/internal/auth/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/cookies"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/claimskeys"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Session returns a snapshot of the authenticated session.
//
// @Summary      Get current session
// @Description  Returns the authenticated user's session snapshot (user identity + roles). Intended for UIs.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  dto.SessionResponse
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /auth/session [get].
func (h *Handler) Session(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerAuthHandler).Start(r.Context(), SpanSessionHandler)
	defer span.End()

	accessToken, err := extractSessionAccessToken(r)
	if err != nil {
		span.RecordError(err)
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}
	if accessToken == "" {
		err := errors.New(ErrMissingAccessToken)
		span.RecordError(err)
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}

	span.SetAttributes(attribute.Bool(AttrAccessTokenPresent, true))

	userID, claims, err := h.Service.Validate(ctx, accessToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrSession)
		httpresponse.WriteAuthError(w, sharederrors.ErrUnauthorized(err.Error()), h.Logger)
		return
	}

	resp := dto.SessionResponse{
		Authenticated: true,
		UserID:        userID,
		Username:      extractStringClaim(claims, claimskeys.Username),
		Email:         extractStringClaim(claims, claimskeys.Email),
		Name:          extractStringClaim(claims, claimskeys.Name),
		Roles:         extractRolesFromClaims(claims),
		ExpiresAt:     extractExp(claims),
	}

	span.SetAttributes(attribute.String(commonkeys.UserID, formatUint(userID)))
	span.AddEvent(EventSessionSuccess)
	span.SetStatus(codes.Ok, StatusSessionSuccess)

	httpresponse.WriteSuccess(w, http.StatusOK, resp, MsgSessionSuccess)
}

func extractSessionAccessToken(r *http.Request) (string, error) {
	// Prefer explicit Authorization header when present.
	if ah := strings.TrimSpace(r.Header.Get("Authorization")); ah != "" {
		parts := strings.SplitN(ah, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && strings.TrimSpace(parts[1]) != "" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	// Backward-compatible cookie support for browser flows.
	return cookies.ExtractAuthToken(r)
}

func extractStringClaim(claims map[string]any, key string) string {
	if claims == nil {
		return ""
	}
	if v, ok := claims[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func extractRolesFromClaims(claims map[string]any) []string {
	if claims == nil {
		return nil
	}
	v, ok := claims[claimskeys.Roles]
	if !ok {
		// backward compat: some middlewares use commonkeys.Roles
		v = claims[commonkeys.Roles]
	}
	switch vv := v.(type) {
	case []string:
		return vv
	case []any:
		roles := make([]string, 0, len(vv))
		for _, it := range vv {
			if s, ok := it.(string); ok {
				roles = append(roles, s)
			}
		}
		return roles
	case string:
		if vv == "" {
			return nil
		}
		return []string{vv}
	default:
		return nil
	}
}

func extractExp(claims map[string]any) *time.Time {
	if claims == nil {
		return nil
	}
	v, ok := claims[claimskeys.Exp]
	if !ok {
		return nil
	}
	switch x := v.(type) {
	case float64:
		t := time.Unix(int64(x), 0).UTC()
		return &t
	case int64:
		t := time.Unix(x, 0).UTC()
		return &t
	case int:
		t := time.Unix(int64(x), 0).UTC()
		return &t
	case string:
		if n, err := parseInt64(x); err == nil {
			t := time.Unix(n, 0).UTC()
			return &t
		}
	}
	return nil
}

func parseInt64(s string) (int64, error) {
	var n int64
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, errors.New("invalid int")
		}
		n = n*10 + int64(r-'0')
	}
	return n, nil
}

func formatUint(v uint64) string {
	const digits = "0123456789"
	if v == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		r := v % 10
		buf = append([]byte{digits[r]}, buf...)
		v /= 10
	}
	return string(buf)
}
