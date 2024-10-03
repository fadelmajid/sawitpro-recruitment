package tests

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/mock"
    "sawitpro-recruitment/models"
)

// MockEstateRepository is a mock implementation of the EstateRepository interface.
type MockEstateRepository struct {
    mock.Mock
}

func (m *MockEstateRepository) CreateEstate(estate *models.Estate) error {
    args := m.Called(estate)
    return args.Error(0)
}

func (m *MockEstateRepository) GetEstateByID(id uuid.UUID) (*models.Estate, error) {
    args := m.Called(id)
    return args.Get(0).(*models.Estate), args.Error(1)
}

func (m *MockEstateRepository) GetEstateStats(id uuid.UUID) (int, int, int, int, error) {
    args := m.Called(id)
    return args.Int(0), args.Int(1), args.Int(2), args.Int(3), args.Error(4)
}

// MockTreeRepository is a mock implementation of the TreeRepository interface.
type MockTreeRepository struct {
    mock.Mock
}

func (m *MockTreeRepository) AddTreeToEstate(tree *models.Tree) error {
    args := m.Called(tree)
    return args.Error(0)
}

func (m *MockTreeRepository) GetTreesByEstateID(estateID uuid.UUID) (map[string]int, error) {
    args := m.Called(estateID)
    return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockTreeRepository) GetTreeByCoordinates(estateID uuid.UUID, x, y int) (*models.Tree, error) {
    args := m.Called(estateID, x, y)
    return args.Get(0).(*models.Tree), args.Error(1)
}