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

func TestGetEstateStats(t *testing.T) {
    e := echo.New()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := handlers.NewEstateHandler(mockEstateRepo)

    t.Run("successful retrieval", func(t *testing.T) {
        estateID := uuid.New()
        estate := &models.Estate{
            ID:     estateID,
            Width:  10,
            Length: 10,
        }
        count, max, min, median := 100, 20, 5, 10

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(estate, nil)
        mockEstateRepo.EXPECT().GetEstateStats(estateID).Return(count, max, min, median, nil)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/stats")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.GetEstateStats(c)) {
            assert.Equal(t, http.StatusOK, rec.Code)
            var response map[string]int
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, count, response["count"])
            assert.Equal(t, max, response["max"])
            assert.Equal(t, min, response["min"])
            assert.Equal(t, median, response["median"])
        }
    })

    t.Run("invalid estate ID format", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/estate/invalid-id/stats", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/stats")
        c.SetParamNames("id")
        c.SetParamValues("invalid-id")

        if assert.NoError(t, handler.GetEstateStats(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid estate ID format", response["message"])
        }
    })

    t.Run("estate not found", func(t *testing.T) {
        estateID := uuid.New()

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(nil, nil)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/stats")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.GetEstateStats(c)) {
            assert.Equal(t, http.StatusNotFound, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Estate not found", response["message"])
        }
    })

    t.Run("database error while retrieving estate", func(t *testing.T) {
        estateID := uuid.New()

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(nil, assert.AnError)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/stats")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.GetEstateStats(c)) {
            assert.Equal(t, http.StatusInternalServerError, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Database error while retrieving estate", response["message"])
        }
    })

    t.Run("database error while fetching estate stats", func(t *testing.T) {
        estateID := uuid.New()
        estate := &models.Estate{
            ID:     estateID,
            Width:  10,
            Length: 10,
        }

        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(estate, nil)
        mockEstateRepo.EXPECT().GetEstateStats(estateID).Return(0, 0, 0, 0, assert.AnError)

        req := httptest.NewRequest(http.MethodGet, "/estate/"+estateID.String()+"/stats", nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/stats")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.GetEstateStats(c)) {
            assert.Equal(t, http.StatusInternalServerError, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Database error while fetching estate stats", response["message"])
        }
    })
}