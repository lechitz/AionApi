// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/output/db/user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/lechitz/AionApi/internal/core/domain"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), ctx, user)
}

// GetAllUsers mocks base method.
func (m *MockUserRepository) GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].([]domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserRepositoryMockRecorder) GetAllUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUserRepository)(nil).GetAllUsers), ctx)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepository) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(ctx domain.ContextControl, ID uint64) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, ID)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), ctx, ID)
}

// GetUserByUsername mocks base method.
func (m *MockUserRepository) GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, username)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsername), ctx, username)
}

// SoftDeleteUser mocks base method.
func (m *MockUserRepository) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDeleteUser indicates an expected call of SoftDeleteUser.
func (mr *MockUserRepositoryMockRecorder) SoftDeleteUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDeleteUser", reflect.TypeOf((*MockUserRepository)(nil).SoftDeleteUser), ctx, userID)
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(ctx domain.ContextControl, userID uint64, fields map[string]interface{}) (domain.UserDomain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, userID, fields)
	ret0, _ := ret[0].(domain.UserDomain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(ctx, userID, fields interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), ctx, userID, fields)
}
