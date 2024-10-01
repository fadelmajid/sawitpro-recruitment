package tests

import (
    "database/sql" // Import sql for database-related operations
    "testing"
	"bytes"        // For handling request body data
    "encoding/json" // For JSON encoding/decoding
    "net/http"     // For HTTP-related functions
    "net/http/httptest" // For creating mock HTTP requests

    "github.com/golang/mock/gomock"
    "github.com/labstack/echo/v4" // Import the Echo web framework
    "github.com/stretchr/testify/assert"
    "sawitpro-recruitment/handlers" // Import your handlers
    "sawitpro-recruitment/models"   // Import your models
    "sawitpro-recruitment/tests/mocks" // Use the generated mock repositories
)

// Test cases for creating an estate
func TestCreateEstate_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    e := echo.New()
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl) // Use the mock repository
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
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
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
    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    treeHandler := handlers.NewTreeHandler(mockTreeRepo)

    // Set up the request payload
    req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": 3, "y": 2, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Expect the GetTreeByCoordinates to be called with estateID = "1234", x = 3, y = 2
    mockTreeRepo.EXPECT().GetTreeByCoordinates("1234", 3, 2).Return(nil, nil).Times(1)

    // Expect AddTreeToEstate to be called
    mockTreeRepo.EXPECT().AddTreeToEstate(gomock.Any()).Return(nil).Times(1)

    // Call the handler and check the result
    if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
    }
}
func TestAddTreeToEstate_TreeAlreadyExists(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    e := echo.New()
    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    treeHandler := handlers.NewTreeHandler(mockTreeRepo)

    req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": 3, "y": 2, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Expect a call to GetTreeByCoordinates that returns an existing tree
    mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 3, 2).Return(&models.Tree{}, nil).Times(1)

    if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestAddTreeToEstate_InvalidInput(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    e := echo.New()
    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    treeHandler := handlers.NewTreeHandler(mockTreeRepo)

    req := httptest.NewRequest(http.MethodPost, "/estate/1234/tree", bytes.NewBufferString(`{"x": -1, "y": 2, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, treeHandler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

// Test case for estate stats success
func TestGetEstateStats_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    e := echo.New()
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    estateHandler := handlers.NewEstateHandler(mockEstateRepo)

    req := httptest.NewRequest(http.MethodGet, "/estate/1234/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Add expectation for GetEstateStats
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
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    estateHandler := handlers.NewEstateHandler(mockEstateRepo)

    req := httptest.NewRequest(http.MethodGet, "/estate/1234/stats", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // Add expectation for GetEstateStats returning sql.ErrNoRows
    mockEstateRepo.EXPECT().GetEstateStats(gomock.Any()).Return(0, 0, 0, 0, sql.ErrNoRows).Times(1)

    if assert.NoError(t, estateHandler.GetEstateStats(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }
}
