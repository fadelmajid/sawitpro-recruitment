package tests

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "sawitpro-recruitment/handlers"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestEstateHandler_CreateEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockEstateRepository(ctrl)
    handler := handlers.NewEstateHandler(mockRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estates", strings.NewReader(`{"width": 100, "length": 200}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    mockRepo.EXPECT().CreateEstate(gomock.Any()).Return(nil)

    if assert.NoError(t, handler.CreateEstate(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
    }
}

func TestEstateHandler_GetEstateStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockEstateRepository(ctrl)
    handler := handlers.NewEstateHandler(mockRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/estates/:id/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(uuid.New().String())

    mockRepo.EXPECT().GetEstateStats(gomock.Any()).Return(10, 50, 5, 25, nil)

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }
}