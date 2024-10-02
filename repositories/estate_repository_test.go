package repositories

import (
    "testing"
    "sawitpro-recruitment/models"
    "sawitpro-recruitment/mocks"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestEstateRepository_CreateEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEstateRepository(ctrl)

    estate := &models.Estate{
        ID:     uuid.New(),
        Width:  100,
        Length: 200,
    }

    mockRepo.EXPECT().CreateEstate(estate).Return(nil)

    err := mockRepo.CreateEstate(estate)
    assert.NoError(t, err)
}

func TestEstateRepository_GetEstateByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEstateRepository(ctrl)

    estateID := uuid.New()
    expectedEstate := &models.Estate{
        ID:     estateID,
        Width:  100,
        Length: 200,
    }

    mockRepo.EXPECT().GetEstateByID(estateID).Return(expectedEstate, nil)

    estate, err := mockRepo.GetEstateByID(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedEstate, estate)
}

func TestEstateRepository_GetEstateStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockEstateRepository(ctrl)

    estateID := uuid.New()
    expectedCount, expectedMax, expectedMin, expectedMedian := 10, 50, 5, 25

    mockRepo.EXPECT().GetEstateStats(estateID).Return(expectedCount, expectedMax, expectedMin, expectedMedian, nil)

    count, max, min, median, err := mockRepo.GetEstateStats(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedCount, count)
    assert.Equal(t, expectedMax, max)
    assert.Equal(t, expectedMin, min)
    assert.Equal(t, expectedMedian, median)
}