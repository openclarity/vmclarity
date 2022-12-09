// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openclarity/vmclarity/backend/pkg/database (interfaces: TargetsTable)

// Package database is a generated GoMock package.
package database

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/openclarity/vmclarity/api/models"
	reflect "reflect"
)

// MockTargetsTable is a mock of TargetsTable interface
type MockTargetsTable struct {
	ctrl     *gomock.Controller
	recorder *MockTargetsTableMockRecorder
}

// MockTargetsTableMockRecorder is the mock recorder for MockTargetsTable
type MockTargetsTableMockRecorder struct {
	mock *MockTargetsTable
}

// NewMockTargetsTable creates a new mock instance
func NewMockTargetsTable(ctrl *gomock.Controller) *MockTargetsTable {
	mock := &MockTargetsTable{ctrl: ctrl}
	mock.recorder = &MockTargetsTableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTargetsTable) EXPECT() *MockTargetsTableMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockTargetsTable) Create(arg0 *Target) (*models.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*models.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockTargetsTableMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTargetsTable)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockTargetsTable) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTargetsTableMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTargetsTable)(nil).Delete), arg0)
}

// Get mocks base method
func (m *MockTargetsTable) Get(arg0 string) (*models.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*models.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockTargetsTableMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTargetsTable)(nil).Get), arg0)
}

// List mocks base method
func (m *MockTargetsTable) List(arg0 models.GetTargetsParams) (*[]models.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*[]models.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockTargetsTableMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTargetsTable)(nil).List), arg0)
}

// Update mocks base method
func (m *MockTargetsTable) Update(arg0 *Target, arg1 string) (*models.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockTargetsTableMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTargetsTable)(nil).Update), arg0, arg1)
}
