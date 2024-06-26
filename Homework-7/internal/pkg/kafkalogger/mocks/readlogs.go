// Code generated by MockGen. DO NOT EDIT.
// Source: ./readlogs.go

// Package mock_logger is a generated GoMock package.
package mock_logger

import (
	reflect "reflect"

	sarama "github.com/IBM/sarama"
	gomock "github.com/golang/mock/gomock"
)

// MockconsumerOps is a mock of consumerOps interface.
type MockconsumerOps struct {
	ctrl     *gomock.Controller
	recorder *MockconsumerOpsMockRecorder
}

// MockconsumerOpsMockRecorder is the mock recorder for MockconsumerOps.
type MockconsumerOpsMockRecorder struct {
	mock *MockconsumerOps
}

// NewMockconsumerOps creates a new mock instance.
func NewMockconsumerOps(ctrl *gomock.Controller) *MockconsumerOps {
	mock := &MockconsumerOps{ctrl: ctrl}
	mock.recorder = &MockconsumerOpsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockconsumerOps) EXPECT() *MockconsumerOpsMockRecorder {
	return m.recorder
}

// ConsumePartition mocks base method.
func (m *MockconsumerOps) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumePartition", topic, partition, offset)
	ret0, _ := ret[0].(sarama.PartitionConsumer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConsumePartition indicates an expected call of ConsumePartition.
func (mr *MockconsumerOpsMockRecorder) ConsumePartition(topic, partition, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumePartition", reflect.TypeOf((*MockconsumerOps)(nil).ConsumePartition), topic, partition, offset)
}

// Partitions mocks base method.
func (m *MockconsumerOps) Partitions(topic string) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Partitions", topic)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Partitions indicates an expected call of Partitions.
func (mr *MockconsumerOpsMockRecorder) Partitions(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Partitions", reflect.TypeOf((*MockconsumerOps)(nil).Partitions), topic)
}
