package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"your-module-name/handlers"
	"your-module-name/models"
	"your-module-name/repositories"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Test cases for creating an estate
func TestCreateEstate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(`{"width": 10, "length": 10}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockEstateRepo.EXPECT().CreateEstate(gomock.Any()).Return(nil).Times(1)

	if assert.NoError(t, estateHandler.CreateEstate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response models.Estate
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Equal(t, 10, response.Width)
		assert.Equal(t, 10, response.Length)
	}
}

func TestCreateEstate_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(`{"width": -5, "length": 10}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, estateHandler.CreateEstate(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// Test cases for adding a tree to an estate
func TestAddTreeToEstate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockTreeRepo := NewMockTreeRepository(ctrl)
	treeHandler := handlers.NewTreeHandler(mockTreeRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": 3, "y": 2, "height": 15}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 3, 2).Return(nil, nil).Times(1)
	mockTreeRepo.EXPECT().AddTreeToEstate(gomock.Any()).Return(nil).Times(1)

	if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestAddTreeToEstate_TreeAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockTreeRepo := NewMockTreeRepository(ctrl)
	treeHandler := handlers.NewTreeHandler(mockTreeRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": 3, "y": 2, "height": 15}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulating that a tree already exists at the specified coordinates
	mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 3, 2).Return(&models.Tree{}, nil).Times(1)

	if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestAddTreeToEstate_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockTreeRepo := NewMockTreeRepository(ctrl)
	treeHandler := handlers.NewTreeHandler(mockTreeRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": -1, "y": 2, "height": 15}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// Test cases for getting estate stats
func TestGetEstateStats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodGet, "/estate/1234/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockEstateRepo.EXPECT().GetEstateStats(gomock.Any()).Return(10, 20, 5, 15, nil).Times(1)

	if assert.NoError(t, estateHandler.GetEstateStats(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]int
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Equal(t, 10, response["count"])
		assert.Equal(t, 20, response["max"])
		assert.Equal(t, 5, response["min"])
		assert.Equal(t, 15, response["median"])
	}
}

func TestGetEstateStats_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodGet, "/estate/1234/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulate that the estate was not found
	mockEstateRepo.EXPECT().GetEstateStats(gomock.Any()).Return(0, 0, 0, 0, sql.ErrNoRows).Times(1)

	if assert.NoError(t, estateHandler.GetEstateStats(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
