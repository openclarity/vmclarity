// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openclarity/vmclarity/backend/pkg/database (interfaces: Database)

// Package database is a generated GoMock package.
package database

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDatabase is a mock of Database interface
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// ScanResultsTable mocks base method
func (m *MockDatabase) ScanResultsTable() ScanResultsTable {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScanResultsTable")
	ret0, _ := ret[0].(ScanResultsTable)
	return ret0
}

// ScanResultsTable indicates an expected call of ScanResultsTable
func (mr *MockDatabaseMockRecorder) ScanResultsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScanResultsTable", reflect.TypeOf((*MockDatabase)(nil).ScanResultsTable))
}

// TargetTable mocks base method
func (m *MockDatabase) TargetTable() TargetTable {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetTable")
	ret0, _ := ret[0].(TargetTable)
	return ret0
}

// TargetTable indicates an expected call of TargetTable
func (mr *MockDatabaseMockRecorder) TargetTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetTable", reflect.TypeOf((*MockDatabase)(nil).TargetTable))
}
