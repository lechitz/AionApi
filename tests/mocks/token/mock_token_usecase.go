// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/usecase/token/token_usecase.go

// Package tokenmocks is a generated GoMock package.
package tokenmocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/lechitz/AionApi/internal/core/domain"
)

// MockTokenUsecase is a mock of TokenUsecase interface.
type MockTokenUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockTokenUsecaseMockRecorder
}

// MockTokenUsecaseMockRecorder is the mock recorder for MockTokenUsecase.
type MockTokenUsecaseMockRecorder struct {
	mock *MockTokenUsecase
}

// NewMockTokenUsecase creates a new mock instance.
func NewMockTokenUsecase(ctrl *gomock.Controller) *MockTokenUsecase {
	mock := &MockTokenUsecase{ctrl: ctrl}
	mock.recorder = &MockTokenUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenUsecase) EXPECT() *MockTokenUsecaseMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockTokenUsecase) CreateToken(ctx context.Context, token domain.TokenDomain) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockTokenUsecaseMockRecorder) CreateToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockTokenUsecase)(nil).CreateToken), ctx, token)
}

// Delete mocks base method.
func (m *MockTokenUsecase) Delete(ctx context.Context, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenUsecaseMockRecorder) Delete(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenUsecase)(nil).Delete), ctx, token)
}

// Save mocks base method.
func (m *MockTokenUsecase) Save(ctx context.Context, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTokenUsecaseMockRecorder) Save(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTokenUsecase)(nil).Save), ctx, token)
}

// Update mocks base method.
func (m *MockTokenUsecase) Update(ctx context.Context, token domain.TokenDomain) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTokenUsecaseMockRecorder) Update(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTokenUsecase)(nil).Update), ctx, token)
}

// VerifyToken mocks base method.
func (m *MockTokenUsecase) VerifyToken(ctx context.Context, token string) (uint64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", ctx, token)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockTokenUsecaseMockRecorder) VerifyToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockTokenUsecase)(nil).VerifyToken), ctx, token)
}
