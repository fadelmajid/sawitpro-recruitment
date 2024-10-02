package tests

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "sawitpro-recruitment/handlers"
    "sawitpro-recruitment/models"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestEstateHandler_GetEstateStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := NewMockEstateRepository(ctrl)
    handler := handlers.NewEstateHandler(mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    // Set up expected calls
    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID)}, nil)
    mockEstateRepo.EXPECT().GetEstateStats(gomock.Any()).Return(10, 20, 5, 15, nil)

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        var response map[string]int
        if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
            assert.Equal(t, 10, response["count"])
            assert.Equal(t, 20, response["max"])
            assert.Equal(t, 5, response["min"])
            assert.Equal(t, 15, response["median"])
        }
    }
}

func TestEstateHandler_GetEstateStats_EstateNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := NewMockEstateRepository(ctrl)
    handler := handlers.NewEstateHandler(mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    // Set up expected calls
    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, nil)

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }
}