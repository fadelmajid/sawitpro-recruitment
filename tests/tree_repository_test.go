package tests

import (
    "testing"
    "sawitpro-recruitment/models"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestTreeRepository_AddTreeToEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockTreeRepository(ctrl)

    tree := &models.Tree{
        ID:       uuid.New(),
        EstateID: uuid.New(),
        X:        10,
        Y:        20,
        Height:   30,
    }

    mockRepo.EXPECT().AddTreeToEstate(tree).Return(nil)

    err := mockRepo.AddTreeToEstate(tree)
    assert.NoError(t, err)
}

func TestTreeRepository_GetTreeByCoordinates(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockTreeRepository(ctrl)

    estateID := uuid.New()
    expectedTree := &models.Tree{
        ID:     uuid.New(),
        X:      10,
        Y:      20,
        Height: 30,
    }

    mockRepo.EXPECT().GetTreeByCoordinates(estateID, 10, 20).Return(expectedTree, nil)

    tree, err := mockRepo.GetTreeByCoordinates(estateID, 10, 20)
    assert.NoError(t, err)
    assert.Equal(t, expectedTree, tree)
}

func TestTreeRepository_GetTreesByEstateID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockTreeRepository(ctrl)

    estateID := uuid.New()
    expectedTreeHeights := map[string]int{
        "10,20": 30,
        "15,25": 35,
    }

    mockRepo.EXPECT().GetTreesByEstateID(estateID).Return(expectedTreeHeights, nil)

    treeHeights, err := mockRepo.GetTreesByEstateID(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedTreeHeights, treeHeights)
}