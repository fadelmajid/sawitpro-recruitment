_package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sawitpro-recruitment/handlers"
	"sawitpro-recruitment/models"
	"sawitpro-recruitment/repositories"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Mocking the repositories
type MockEstateRepository struct {
	repositories.EstateRepository
}

type MockTreeRepository struct {
	repositories.TreeRepository
}

func TestCreateEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	// Mock Estate Repository
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewBufferString(`{"width": 10, "length": 10}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Expect the repository to create an estate successfully
	mockEstateRepo.EXPECT().CreateEstate(gomock.Any()).Return(nil).Times(1)

	if assert.NoError(t, estateHandler.CreateEstate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response models.Estate
		json.NewDecoder(rec.Body).Decode(&response)
		assert.Equal(t, 10, response.Width)
		assert.Equal(t, 10, response.Length)
	}
}

func TestAddTreeToEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	// Mock Tree Repository
	mockTreeRepo := NewMockTreeRepository(ctrl)
	treeHandler := handlers.NewTreeHandler(mockTreeRepo)

	req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": 3, "y": 2, "height": 15}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Expect the repository to check if the tree exists and then add the tree
	mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 3, 2).Return(nil, nil).Times(1)
	mockTreeRepo.EXPECT().AddTreeToEstate(gomock.Any()).Return(nil).Times(1)

	if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestGetEstateStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()

	// Mock Estate Repository
	mockEstateRepo := NewMockEstateRepository(ctrl)
	estateHandler := handlers.NewEstateHandler(mockEstateRepo)

	req := httptest.NewRequest(http.MethodGet, "/estate/1234/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Expect the repository to return some stats
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
