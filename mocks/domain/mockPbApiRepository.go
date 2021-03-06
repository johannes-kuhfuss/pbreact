// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/johannes-kuhfuss/pbreact/domain (interfaces: PbApiRepository)

// Package domain is a generated GoMock package.
package domain

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/johannes-kuhfuss/pbreact/dto"
	api_error "github.com/johannes-kuhfuss/services_utils/api_error"
)

// MockPbApiRepository is a mock of PbApiRepository interface.
type MockPbApiRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPbApiRepositoryMockRecorder
}

// MockPbApiRepositoryMockRecorder is the mock recorder for MockPbApiRepository.
type MockPbApiRepositoryMockRecorder struct {
	mock *MockPbApiRepository
}

// NewMockPbApiRepository creates a new mock instance.
func NewMockPbApiRepository(ctrl *gomock.Controller) *MockPbApiRepository {
	mock := &MockPbApiRepository{ctrl: ctrl}
	mock.recorder = &MockPbApiRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPbApiRepository) EXPECT() *MockPbApiRepositoryMockRecorder {
	return m.recorder
}

// GetNotifications mocks base method.
func (m *MockPbApiRepository) GetNotifications() (*dto.PbSubscriptionResponse, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotifications")
	ret0, _ := ret[0].(*dto.PbSubscriptionResponse)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// GetNotifications indicates an expected call of GetNotifications.
func (mr *MockPbApiRepositoryMockRecorder) GetNotifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotifications", reflect.TypeOf((*MockPbApiRepository)(nil).GetNotifications))
}

// RegisterForNotifications mocks base method.
func (m *MockPbApiRepository) RegisterForNotifications() api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterForNotifications")
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// RegisterForNotifications indicates an expected call of RegisterForNotifications.
func (mr *MockPbApiRepositoryMockRecorder) RegisterForNotifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterForNotifications", reflect.TypeOf((*MockPbApiRepository)(nil).RegisterForNotifications))
}

// UnregisterForNotifications mocks base method.
func (m *MockPbApiRepository) UnregisterForNotifications(arg0 dto.PbSubscriptionResponse) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterForNotifications", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// UnregisterForNotifications indicates an expected call of UnregisterForNotifications.
func (mr *MockPbApiRepositoryMockRecorder) UnregisterForNotifications(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterForNotifications", reflect.TypeOf((*MockPbApiRepository)(nil).UnregisterForNotifications), arg0)
}
