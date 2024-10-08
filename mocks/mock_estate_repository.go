// Code generated by MockGen. DO NOT EDIT.
// Source: repositories/estate_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	models "sawitpro-recruitment/models"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockEstateRepository is a mock of EstateRepository interface.
type MockEstateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEstateRepositoryMockRecorder
}

// MockEstateRepositoryMockRecorder is the mock recorder for MockEstateRepository.
type MockEstateRepositoryMockRecorder struct {
	mock *MockEstateRepository
}

// NewMockEstateRepository creates a new mock instance.
func NewMockEstateRepository(ctrl *gomock.Controller) *MockEstateRepository {
	mock := &MockEstateRepository{ctrl: ctrl}
	mock.recorder = &MockEstateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEstateRepository) EXPECT() *MockEstateRepositoryMockRecorder {
	return m.recorder
}

// CreateEstate mocks base method.
func (m *MockEstateRepository) CreateEstate(estate *models.Estate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEstate", estate)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEstate indicates an expected call of CreateEstate.
func (mr *MockEstateRepositoryMockRecorder) CreateEstate(estate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEstate", reflect.TypeOf((*MockEstateRepository)(nil).CreateEstate), estate)
}

// GetEstateByID mocks base method.
func (m *MockEstateRepository) GetEstateByID(id uuid.UUID) (*models.Estate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEstateByID", id)
	ret0, _ := ret[0].(*models.Estate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEstateByID indicates an expected call of GetEstateByID.
func (mr *MockEstateRepositoryMockRecorder) GetEstateByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEstateByID", reflect.TypeOf((*MockEstateRepository)(nil).GetEstateByID), id)
}

// GetEstateStats mocks base method.
func (m *MockEstateRepository) GetEstateStats(id uuid.UUID) (int, int, int, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEstateStats", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(int)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetEstateStats indicates an expected call of GetEstateStats.
func (mr *MockEstateRepositoryMockRecorder) GetEstateStats(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEstateStats", reflect.TypeOf((*MockEstateRepository)(nil).GetEstateStats), id)
}
