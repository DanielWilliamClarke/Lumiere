// Code generated by MockGen. DO NOT EDIT.
// Source: ./utils/generator.go

// Package mock_utils is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockICodeGenerator is a mock of ICodeGenerator interface
type MockICodeGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockICodeGeneratorMockRecorder
}

// MockICodeGeneratorMockRecorder is the mock recorder for MockICodeGenerator
type MockICodeGeneratorMockRecorder struct {
	mock *MockICodeGenerator
}

// NewMockICodeGenerator creates a new mock instance
func NewMockICodeGenerator(ctrl *gomock.Controller) *MockICodeGenerator {
	mock := &MockICodeGenerator{ctrl: ctrl}
	mock.recorder = &MockICodeGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockICodeGenerator) EXPECT() *MockICodeGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockICodeGenerator) Generate(length int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", length)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (mr *MockICodeGeneratorMockRecorder) Generate(length interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockICodeGenerator)(nil).Generate), length)
}
