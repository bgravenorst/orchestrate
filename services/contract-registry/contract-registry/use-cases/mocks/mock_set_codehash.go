// Code generated by MockGen. DO NOT EDIT.
// Source: set_codehash.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	reflect "reflect"
)

// MockSetCodeHashUseCase is a mock of SetCodeHashUseCase interface
type MockSetCodeHashUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSetCodeHashUseCaseMockRecorder
}

// MockSetCodeHashUseCaseMockRecorder is the mock recorder for MockSetCodeHashUseCase
type MockSetCodeHashUseCaseMockRecorder struct {
	mock *MockSetCodeHashUseCase
}

// NewMockSetCodeHashUseCase creates a new mock instance
func NewMockSetCodeHashUseCase(ctrl *gomock.Controller) *MockSetCodeHashUseCase {
	mock := &MockSetCodeHashUseCase{ctrl: ctrl}
	mock.recorder = &MockSetCodeHashUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSetCodeHashUseCase) EXPECT() *MockSetCodeHashUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockSetCodeHashUseCase) Execute(ctx context.Context, account *common.AccountInstance, hash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, account, hash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockSetCodeHashUseCaseMockRecorder) Execute(ctx, account, hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSetCodeHashUseCase)(nil).Execute), ctx, account, hash)
}