package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "sawitpro-recruitment/handlers"
    "sawitpro-recruitment/models"
)

// MockEstateRepository is a mock implementation of the EstateRepository interface.
type MockEstateRepository struct {
    mock.Mock
}

func (m *MockEstateRepository) CreateEstate(estate *models.Estate) error {
    args := m.Called(estate)
    return args.Error(0)
}

func (m *MockEstateRepository) GetEstateByID(id uuid.UUID) (*models.Estate, error) {
    args := m.Called(id)
    return args.Get(0).(*models.Estate), args.Error(1)
}

func (m *MockEstateRepository) GetEstateStats(id uuid.UUID) (int, int, int, int, error) {
    args := m.Called(id)
    return args.Int(0), args.Int(1), args.Int(2), args.Int(3), args.Error(4)
}

func TestCreateEstate(t *testing.T) {
    e := echo.New()
    mockRepo := new(MockEstateRepository)
    handler := handlers.NewEstateHandler(mockRepo)

    t.Run("successful creation", func(t *testing.T) {
        estate := &models.Estate{
            Width:  100,
            Length: 100,
        }
        mockRepo.On("CreateEstate", mock.AnythingOfType("*models.Estate")).Return(nil)

        body, _ := json.Marshal(estate)
        req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        if assert.NoError(t, handler.CreateEstate(c)) {
            assert.Equal(t, http.StatusCreated, rec.Code)
            var response models.Estate
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, estate.Width, response.Width)
            assert.Equal(t, estate.Length, response.Length)
            assert.NotEmpty(t, response.ID)
        }

        mockRepo.AssertExpectations(t)
    })

    t.Run("invalid input format", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader([]byte("invalid")))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        if assert.NoError(t, handler.CreateEstate(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid input format", response["message"])
        }
    })

    t.Run("invalid estate dimensions", func(t *testing.T) {
        estate := &models.Estate{
            Width:  0,
            Length: 100,
        }
        body, _ := json.Marshal(estate)
        req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)

        if assert.NoError(t, handler.CreateEstate(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Estate dimensions must be between 1 and 50000", response["message"])
        }
    })
}