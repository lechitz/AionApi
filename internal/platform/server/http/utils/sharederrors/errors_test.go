package sharederrors_test

import (
	"errors"
	"net/http"
	"testing"

	httperrors "github.com/lechitz/AionApi/internal/platform/server/http/errors"
	sharederrors "github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
)

func TestErrorTypesAndConstructors(t *testing.T) {
	if sharederrors.ErrMissingUserID().Error() != sharederrors.ErrMsgMissingUserID {
		t.Fatalf("unexpected missing user id error message")
	}

	if sharederrors.ErrUnauthorized("bad token").Error() != "unauthorized: bad token" {
		t.Fatalf("unexpected unauthorized error message")
	}

	if sharederrors.ErrForbidden("role").Error() != "forbidden: role" {
		t.Fatalf("unexpected forbidden error message")
	}

	if sharederrors.NewAuthenticationError("expired").Error() != "unauthorized: expired" {
		t.Fatalf("unexpected authentication error message")
	}

	if sharederrors.NewValidationError("field", "bad").Error() != "validation error on field: bad" {
		t.Fatalf("unexpected validation error message")
	}
}

func TestValidationErrorFormats(t *testing.T) {
	if (&sharederrors.ValidationError{Field: "name"}).Error() != "validation error on name" {
		t.Fatalf("unexpected validation format for field-only")
	}

	if (&sharederrors.ValidationError{}).Error() != "validation error" {
		t.Fatalf("unexpected validation format for empty")
	}
}

func TestMissingFieldsHelpers(t *testing.T) {
	err := sharederrors.MissingFields("username")
	if err == nil || err.Error() != "validation error on username: missing field: username" {
		t.Fatalf("unexpected single missing field error: %v", err)
	}

	err = sharederrors.MissingFields("username", "email")
	if err == nil || err.Error() != "validation error on fields: required fields missing: username, email" {
		t.Fatalf("unexpected multi missing field error: %v", err)
	}

	err = sharederrors.AtLeastOneFieldRequired("username", "email")
	if err == nil || err.Error() != "at least one of the following fields must be provided: username, email" {
		t.Fatalf("unexpected at-least-one error: %v", err)
	}
}

func TestMapErrorToHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: http.StatusOK},
		{name: "validation type", err: &sharederrors.ValidationError{}, want: http.StatusBadRequest},
		{name: "unauthorized type", err: &sharederrors.UnauthorizedError{}, want: http.StatusUnauthorized},
		{name: "forbidden type", err: &sharederrors.ForbiddenError{}, want: http.StatusForbidden},
		{name: "missing user id type", err: &sharederrors.MissingUserIDError{}, want: http.StatusUnauthorized},
		{name: "auth type", err: &sharederrors.AuthenticationError{}, want: http.StatusUnauthorized},
		{name: "not found sentinel", err: httperrors.ErrResourceNotFound, want: http.StatusNotFound},
		{name: "method sentinel", err: httperrors.ErrMethodNotAllowed, want: http.StatusMethodNotAllowed},
		{name: "parse user id sentinel", err: sharederrors.ErrParseUserID, want: http.StatusBadRequest},
		{name: "no fields sentinel", err: sharederrors.ErrNoFieldsToUpdate, want: http.StatusBadRequest},
		{name: "username exists", err: sharederrors.ErrUsernameExists, want: http.StatusConflict},
		{name: "email exists", err: sharederrors.ErrEmailExists, want: http.StatusConflict},
		{name: "domain conflict", err: sharederrors.ErrDomainConflict, want: http.StatusConflict},
		{name: "unknown", err: errors.New("x"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sharederrors.MapErrorToHTTPStatus(tt.err)
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
