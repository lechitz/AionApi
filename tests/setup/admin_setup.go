// Package setup provides test suite builders and common test helpers for unit tests.
package setup

import (
	"context"
	"testing"
	"time"

	admindomain "github.com/lechitz/aion-api/internal/admin/core/domain"
	"github.com/lechitz/aion-api/internal/admin/core/usecase"
	"github.com/lechitz/aion-api/tests/mocks"
	"go.uber.org/mock/gomock"
)

// AdminServiceTestSuite groups mocked dependencies and the system under test (AdminService)
// to keep admin-related tests concise and consistent.
type AdminServiceTestSuite struct {
	Ctrl            *gomock.Controller
	Logger          *mocks.MockContextLogger
	AdminRepository *mocks.MockAdminRepository
	RoleCache       *mocks.MockRoleCache
	SessionRevoker  *mocks.MockSessionRevoker
	AdminService    *usecase.Service
	Ctx             context.Context
}

// AdminServiceTest initializes and returns an AdminServiceTestSuite using mocked output ports.
// Use this helper to bootstrap each test and ensure proper teardown via Ctrl.Finish().
func AdminServiceTest(t *testing.T) *AdminServiceTestSuite {
	ctrl := gomock.NewController(t)

	adminRepo := mocks.NewMockAdminRepository(ctrl)
	roleCache := mocks.NewMockRoleCache(ctrl)
	sessionRevoker := mocks.NewMockSessionRevoker(ctrl)
	log := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger (no-ops).
	ExpectLoggerDefaultBehavior(log)

	svc := usecase.NewService(
		adminRepo,
		roleCache,
		sessionRevoker,
		log,
	)

	return &AdminServiceTestSuite{
		Ctrl:            ctrl,
		Logger:          log,
		AdminRepository: adminRepo,
		RoleCache:       roleCache,
		SessionRevoker:  sessionRevoker,
		AdminService:    svc,
		Ctx:             t.Context(),
	}
}

// DefaultTestUserWithRoles returns a valid domain.AdminUser with roles commonly used in admin tests.
func DefaultTestUserWithRoles(roles []string) admindomain.AdminUser {
	return admindomain.AdminUser{
		ID:        1,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "user@example.com",
		Roles:     roles,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
