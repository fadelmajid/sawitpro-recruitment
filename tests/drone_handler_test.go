package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "sawitpro-recruitment/handlers"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestCalculateDronePlanWithLimit(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockTreeRepository(ctrl)
    handler := handlers.NewDroneHandler(mockRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/estates/:id/drone-plan", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(uuid.New().String())

    mockRepo.EXPECT().GetTreesByEstateID(gomock.Any()).Return(map[string]int{
        "1,1": 10,
        "1,2": 15,
    }, nil)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }
}