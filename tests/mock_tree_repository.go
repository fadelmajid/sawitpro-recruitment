// Code generated by MockGen. DO NOT EDIT.
// Source: repositories/tree_repository.go

// Package tests is a generated GoMock package.
package tests

import (
    reflect "reflect"
    models "sawitpro-recruitment/models"
    uuid "github.com/google/uuid"
    gomock "github.com/golang/mock/gomock"
)

// MockTreeRepository is a mock of TreeRepository interface.
type MockTreeRepository struct {
    ctrl     *gomock.Controller
    recorder *MockTreeRepositoryMockRecorder
}

// MockTreeRepositoryMockRecorder is the mock recorder for MockTreeRepository.
type MockTreeRepositoryMockRecorder struct {
    mock *MockTreeRepository
}

// NewMockTreeRepository creates a new mock instance.
func NewMockTreeRepository(ctrl *gomock.Controller) *MockTreeRepository {
    mock := &MockTreeRepository{ctrl: ctrl}
    mock.recorder = &MockTreeRepositoryMockRecorder{mock}
    return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTreeRepository) EXPECT() *MockTreeRepositoryMockRecorder {
    return m.recorder
}

// AddTreeToEstate mocks base method.
func (m *MockTreeRepository) AddTreeToEstate(tree *models.Tree) error {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "AddTreeToEstate", tree)
    ret0, _ := ret[0].(error)
    return ret0
}

// AddTreeToEstate indicates an expected call of AddTreeToEstate.
func (mr *MockTreeRepositoryMockRecorder) AddTreeToEstate(tree interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTreeToEstate", reflect.TypeOf((*MockTreeRepository)(nil).AddTreeToEstate), tree)
}

// GetTreeByCoordinates mocks base method.
func (m *MockTreeRepository) GetTreeByCoordinates(estateID uuid.UUID, x, y int) (*models.Tree, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "GetTreeByCoordinates", estateID, x, y)
    ret0, _ := ret[0].(*models.Tree)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// GetTreeByCoordinates indicates an expected call of GetTreeByCoordinates.
func (mr *MockTreeRepositoryMockRecorder) GetTreeByCoordinates(estateID, x, y interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTreeByCoordinates", reflect.TypeOf((*MockTreeRepository)(nil).GetTreeByCoordinates), estateID, x, y)
}

// GetTreesByEstateID mocks base method.
func (m *MockTreeRepository) GetTreesByEstateID(estateID uuid.UUID) (map[string]int, error) {
    m.ctrl.T.Helper()
    ret := m.ctrl.Call(m, "GetTreesByEstateID", estateID)
    ret0, _ := ret[0].(map[string]int)
    ret1, _ := ret[1].(error)
    return ret0, ret1
}

// GetTreesByEstateID indicates an expected call of GetTreesByEstateID.
func (mr *MockTreeRepositoryMockRecorder) GetTreesByEstateID(estateID interface{}) *gomock.Call {
    mr.mock.ctrl.T.Helper()
    return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTreesByEstateID", reflect.TypeOf((*MockTreeRepository)(nil).GetTreesByEstateID), estateID)
}