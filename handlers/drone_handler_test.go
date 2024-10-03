package handlers

import (
    "errors"
    "net/http"
    "net/http/httptest"
    "encoding/json"
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

func TestCalculateDronePlanWithLimit_InvalidEstateID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewDroneHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    invalidEstateID := "invalid-uuid"
    req := httptest.NewRequest(http.MethodGet, "/estate/"+invalidEstateID+"/drone-plan", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(invalidEstateID)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestCalculateDronePlanWithLimit_EstateNotFound(t *testing.T) {
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

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, nil)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }
}

func TestCalculateDronePlanWithLimit_DatabaseErrorRetrievingEstate(t *testing.T) {
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

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, errors.New("database error"))

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusInternalServerError, rec.Code)
    }
}

func TestCalculateDronePlanWithLimit_DatabaseErrorFetchingTreeHeights(t *testing.T) {
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

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 50, Length: 50}, nil)
    mockTreeRepo.EXPECT().GetTreesByEstateID(gomock.Any()).Return(nil, errors.New("database error"))

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusInternalServerError, rec.Code)
    }
}

func TestCalculateDronePlanWithLimit_InvalidMaxDistance(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewDroneHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/drone-plan?max_distance=invalid", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestCalculateDronePlanWithLimit_MaxDistanceReached(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewDroneHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID+"/drone-plan?max_distance=40", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 50, Length: 50}, nil)
    mockTreeRepo.EXPECT().GetTreesByEstateID(gomock.Any()).Return(map[string]int{
        "1,1": 10,
        "1,2": 15,
        "2,1": 20,
        "2,2": 25,
        "3,1": 30,
        "3,2": 35,
    }, nil)

    if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        var response map[string]interface{}
        if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
            assert.Equal(t, 40, int(response["total_distance"].(float64)))
            landedAt := response["landed_at"].(map[string]interface{})
            assert.Equal(t, 3, int(landedAt["x"].(float64)))
            assert.Equal(t, 1, int(landedAt["y"].(float64)))
        }
    }
}