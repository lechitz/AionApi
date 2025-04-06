// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/output/cache/token.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/lechitz/AionApi/internal/core/domain"
)

// MockTokenCreator is a mock of TokenCreator interface.
type MockTokenCreator struct {
	ctrl     *gomock.Controller
	recorder *MockTokenCreatorMockRecorder
}

// MockTokenCreatorMockRecorder is the mock recorder for MockTokenCreator.
type MockTokenCreatorMockRecorder struct {
	mock *MockTokenCreator
}

// NewMockTokenCreator creates a new mock instance.
func NewMockTokenCreator(ctrl *gomock.Controller) *MockTokenCreator {
	mock := &MockTokenCreator{ctrl: ctrl}
	mock.recorder = &MockTokenCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenCreator) EXPECT() *MockTokenCreatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTokenCreator) CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTokenCreatorMockRecorder) Create(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockTokenCreator)(nil).CreateToken), ctx, token)
}

// MockTokenChecker is a mock of TokenChecker interface.
type MockTokenChecker struct {
	ctrl     *gomock.Controller
	recorder *MockTokenCheckerMockRecorder
}

// MockTokenCheckerMockRecorder is the mock recorder for MockTokenChecker.
type MockTokenCheckerMockRecorder struct {
	mock *MockTokenChecker
}

// NewMockTokenChecker creates a new mock instance.
func NewMockTokenChecker(ctrl *gomock.Controller) *MockTokenChecker {
	mock := &MockTokenChecker{ctrl: ctrl}
	mock.recorder = &MockTokenCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenChecker) EXPECT() *MockTokenCheckerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockTokenChecker) Check(ctx domain.ContextControl, token string) (uint64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, token)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Check indicates an expected call of Check.
func (mr *MockTokenCheckerMockRecorder) Check(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockTokenChecker)(nil).Check), ctx, token)
}

// MockTokenSaver is a mock of TokenSaver interface.
type MockTokenSaver struct {
	ctrl     *gomock.Controller
	recorder *MockTokenSaverMockRecorder
}

// MockTokenSaverMockRecorder is the mock recorder for MockTokenSaver.
type MockTokenSaverMockRecorder struct {
	mock *MockTokenSaver
}

// NewMockTokenSaver creates a new mock instance.
func NewMockTokenSaver(ctrl *gomock.Controller) *MockTokenSaver {
	mock := &MockTokenSaver{ctrl: ctrl}
	mock.recorder = &MockTokenSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenSaver) EXPECT() *MockTokenSaverMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockTokenSaver) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTokenSaverMockRecorder) Save(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTokenSaver)(nil).Save), ctx, token)
}

// MockTokenUpdater is a mock of TokenUpdater interface.
type MockTokenUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockTokenUpdaterMockRecorder
}

// MockTokenUpdaterMockRecorder is the mock recorder for MockTokenUpdater.
type MockTokenUpdaterMockRecorder struct {
	mock *MockTokenUpdater
}

// NewMockTokenUpdater creates a new mock instance.
func NewMockTokenUpdater(ctrl *gomock.Controller) *MockTokenUpdater {
	mock := &MockTokenUpdater{ctrl: ctrl}
	mock.recorder = &MockTokenUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenUpdater) EXPECT() *MockTokenUpdaterMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockTokenUpdater) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTokenUpdaterMockRecorder) Update(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTokenUpdater)(nil).Update), ctx, token)
}

// MockTokenDeleter is a mock of TokenDeleter interface.
type MockTokenDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockTokenDeleterMockRecorder
}

// MockTokenDeleterMockRecorder is the mock recorder for MockTokenDeleter.
type MockTokenDeleterMockRecorder struct {
	mock *MockTokenDeleter
}

// NewMockTokenDeleter creates a new mock instance.
func NewMockTokenDeleter(ctrl *gomock.Controller) *MockTokenDeleter {
	mock := &MockTokenDeleter{ctrl: ctrl}
	mock.recorder = &MockTokenDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenDeleter) EXPECT() *MockTokenDeleterMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockTokenDeleter) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenDeleterMockRecorder) Delete(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenDeleter)(nil).Delete), ctx, token)
}

// MockTokenService is a mock of TokenService interface.
type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
}

// MockTokenServiceMockRecorder is the mock recorder for MockTokenService.
type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

// NewMockTokenService creates a new mock instance.
func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	mock := &MockTokenService{ctrl: ctrl}
	mock.recorder = &MockTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockTokenService) Check(ctx domain.ContextControl, token string) (uint64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, token)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Check indicates an expected call of Check.
func (mr *MockTokenServiceMockRecorder) Check(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockTokenService)(nil).Check), ctx, token)
}

// Create mocks base method.
func (m *MockTokenService) CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTokenServiceMockRecorder) Create(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockTokenService)(nil).CreateToken), ctx, token)
}

// Delete mocks base method.
func (m *MockTokenService) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenServiceMockRecorder) Delete(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenService)(nil).Delete), ctx, token)
}

// Save mocks base method.
func (m *MockTokenService) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTokenServiceMockRecorder) Save(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTokenService)(nil).Save), ctx, token)
}

// Update mocks base method.
func (m *MockTokenService) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTokenServiceMockRecorder) Update(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTokenService)(nil).Update), ctx, token)
}
