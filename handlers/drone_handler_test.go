package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "sawitpro-recruitment/models"
    "sawitpro-recruitment/mocks"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestCalculateDronePlanWithLimit(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewDroneHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/drone-plan", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    // Set up expected calls
    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 50, Length: 50}, nil)
    mockTreeRepo.EXPECT().GetTreesByEstateID(gomock.Any()).Return(map[string]int{
        "1,1": 10,
        "1,2": 15,
    }, nil)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }
}