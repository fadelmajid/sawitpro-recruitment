package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "sawitpro-recruitment/models"
    "sawitpro-recruitment/mocks"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestEstateHandler_CreateEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 100, "length": 200}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    mockEstateRepo.EXPECT().CreateEstate(gomock.Any()).Return(nil)

    if assert.NoError(t, handler.CreateEstate(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
        var response map[string]string
        if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
            assert.NotEmpty(t, response["Id"])
        }
    }
}

func TestEstateHandler_CreateEstate_InvalidInput(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`invalid json`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, handler.CreateEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestEstateHandler_CreateEstate_InvalidDimensions(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 0, "length": 200}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, handler.CreateEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestEstateHandler_CreateEstate_DatabaseError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(`{"width": 100, "length": 200}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    mockEstateRepo.EXPECT().CreateEstate(gomock.Any()).Return(errors.New("database error"))

    if assert.NoError(t, handler.CreateEstate(c)) {
        assert.Equal(t, http.StatusInternalServerError, rec.Code)
    }
}

func TestEstateHandler_GetEstateStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

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

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, nil)

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }
}

func TestEstateHandler_GetEstateStats_InvalidID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/estate/invalid-id/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("invalid-id")

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestEstateHandler_GetEstateStats_ErrorFetchingStats(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewEstateHandler(mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID)}, nil)
    mockEstateRepo.EXPECT().GetEstateStats(gomock.Any()).Return(0, 0, 0, 0, errors.New("database error"))

    if assert.NoError(t, handler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusInternalServerError, rec.Code)
    }
}