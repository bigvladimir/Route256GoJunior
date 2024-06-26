// Code generated by MockGen. DO NOT EDIT.
// Source: ./server.go

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	dto "homework/internal/app/pvz/dto"
	kafkalogger "homework/internal/pkg/kafkalogger"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockcoreOps is a mock of coreOps interface.
type MockcoreOps struct {
	ctrl     *gomock.Controller
	recorder *MockcoreOpsMockRecorder
}

// MockcoreOpsMockRecorder is the mock recorder for MockcoreOps.
type MockcoreOpsMockRecorder struct {
	mock *MockcoreOps
}

// NewMockcoreOps creates a new mock instance.
func NewMockcoreOps(ctrl *gomock.Controller) *MockcoreOps {
	mock := &MockcoreOps{ctrl: ctrl}
	mock.recorder = &MockcoreOpsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcoreOps) EXPECT() *MockcoreOpsMockRecorder {
	return m.recorder
}

// AddPvz mocks base method.
func (m *MockcoreOps) AddPvz(ctx context.Context, input dto.PvzInput) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPvz", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPvz indicates an expected call of AddPvz.
func (mr *MockcoreOpsMockRecorder) AddPvz(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPvz", reflect.TypeOf((*MockcoreOps)(nil).AddPvz), ctx, input)
}

// DeletePvz mocks base method.
func (m *MockcoreOps) DeletePvz(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePvz", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePvz indicates an expected call of DeletePvz.
func (mr *MockcoreOpsMockRecorder) DeletePvz(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePvz", reflect.TypeOf((*MockcoreOps)(nil).DeletePvz), ctx, id)
}

// GetPvzByID mocks base method.
func (m *MockcoreOps) GetPvzByID(ctx context.Context, id int64) (dto.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvzByID", ctx, id)
	ret0, _ := ret[0].(dto.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvzByID indicates an expected call of GetPvzByID.
func (mr *MockcoreOpsMockRecorder) GetPvzByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvzByID", reflect.TypeOf((*MockcoreOps)(nil).GetPvzByID), ctx, id)
}

// LogMessage mocks base method.
func (m *MockcoreOps) LogMessage(message kafkalogger.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogMessage", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogMessage indicates an expected call of LogMessage.
func (mr *MockcoreOpsMockRecorder) LogMessage(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMessage", reflect.TypeOf((*MockcoreOps)(nil).LogMessage), message)
}

// LogMessages mocks base method.
func (m *MockcoreOps) LogMessages(message []kafkalogger.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogMessages", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogMessages indicates an expected call of LogMessages.
func (mr *MockcoreOpsMockRecorder) LogMessages(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogMessages", reflect.TypeOf((*MockcoreOps)(nil).LogMessages), message)
}

// ModifyPvz mocks base method.
func (m *MockcoreOps) ModifyPvz(ctx context.Context, input dto.Pvz) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifyPvz", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModifyPvz indicates an expected call of ModifyPvz.
func (mr *MockcoreOpsMockRecorder) ModifyPvz(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifyPvz", reflect.TypeOf((*MockcoreOps)(nil).ModifyPvz), ctx, input)
}

// UpdatePvz mocks base method.
func (m *MockcoreOps) UpdatePvz(ctx context.Context, input dto.Pvz) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePvz", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePvz indicates an expected call of UpdatePvz.
func (mr *MockcoreOpsMockRecorder) UpdatePvz(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePvz", reflect.TypeOf((*MockcoreOps)(nil).UpdatePvz), ctx, input)
}
