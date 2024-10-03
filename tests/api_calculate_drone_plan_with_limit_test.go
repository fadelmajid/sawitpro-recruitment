package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"

    "github.com/golang/mock/gomock"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/google/uuid"
    "sawitpro-recruitment/handlers"
    "sawitpro-recruitment/mocks"
    "sawitpro-recruitment/models"
)

func TestCalculateDronePlanWithLimit(t *testing.T) {
    e := echo.New()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := handlers.NewDroneHandler(mockTreeRepo, mockEstateRepo)

    t.Run("successful calculation without limit", func(t *testing.T) {
        estateID := uuid.New()
        estate := &models.Estate{
            ID:     estateID,
            Width:  10,
            Length: 10,
        }
        treeHeights := map[string]int{
            "1,1": 10, "1,2": 15, "1,3": 20,
            // Add more tree heights as needed
        }

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(estate, nil)
        mockTreeRepo.EXPECT().GetTreesByEstateID(estateID).Return(treeHeights, nil)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusOK, rec.Code)
            var response map[string]interface{}
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.NotNil(t, response["distance"])
        }
    })

    t.Run("invalid estate ID format", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/estate/invalid-id/drone-plan", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues("invalid-id")

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid estate ID format", response["message"])
        }
    })

    t.Run("estate not found", func(t *testing.T) {
        estateID := uuid.New()

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(nil, nil)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusNotFound, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Estate not found", response["message"])
        }
    })

    t.Run("database error while retrieving estate", func(t *testing.T) {
        estateID := uuid.New()

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(nil, assert.AnError)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusInternalServerError, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Database error while retrieving estate", response["message"])
        }
    })

    t.Run("database error while fetching tree heights", func(t *testing.T) {
        estateID := uuid.New()
        estate := &models.Estate{
            ID:     estateID,
            Width:  10,
            Length: 10,
        }

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(estate, nil)
        mockTreeRepo.EXPECT().GetTreesByEstateID(estateID).Return(nil, assert.AnError)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusInternalServerError, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Database error while fetching tree heights", response["message"])
        }
    })

    t.Run("invalid max_distance value", func(t *testing.T) {
        estateID := uuid.New()
        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan?max_distance=invalid", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid max_distance value", response["message"])
        }
    })

    t.Run("max distance reached", func(t *testing.T) {
        estateID := uuid.New()
        estate := &models.Estate{
            ID:     estateID,
            Width:  10,
            Length: 10,
        }
        treeHeights := map[string]int{
            "1,1": 10, "1,2": 15, "2,1": 20, "2,2": 25,
            "3,1": 30, "3,2": 35,
        }

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(estate, nil)
        mockTreeRepo.EXPECT().GetTreesByEstateID(estateID).Return(treeHeights, nil)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/drone-plan?max_distance=40", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/drone-plan")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.CalculateDronePlanWithLimit(c)) {
            assert.Equal(t, http.StatusOK, rec.Code)
            var response map[string]interface{}
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, 40, int(response["distance"].(float64)))
            landedAt := response["rest"].(map[string]interface{})
            assert.Equal(t, 3, int(landedAt["x"].(float64)))
            assert.Equal(t, 1, int(landedAt["y"].(float64)))
        }
    })
}