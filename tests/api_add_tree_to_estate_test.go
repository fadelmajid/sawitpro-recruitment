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
    "github.com/golang/mock/gomock"
    "sawitpro-recruitment/handlers"
    "sawitpro-recruitment/models"
    "sawitpro-recruitment/mocks"
)

func TestAddTreeToEstate(t *testing.T) {
    e := echo.New()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := handlers.NewTreeHandler(mockTreeRepo, mockEstateRepo)

    t.Run("successful addition", func(t *testing.T) {
        estateID := uuid.New()
        tree := &models.Tree{
            EstateID: estateID,
            X:        1,
            Y:        1,
            Height:   10,
        }
        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(&models.Estate{ID: estateID, Width: 100, Length: 100}, nil)
        mockTreeRepo.EXPECT().GetTreeByCoordinates(estateID, tree.X, tree.Y).Return(nil, nil)
        mockTreeRepo.EXPECT().AddTreeToEstate(gomock.Any()).DoAndReturn(func(t *models.Tree) error {
            t.ID = uuid.New()
            return nil
        })

        body, _ := json.Marshal(tree)
        req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID.String()+"/tree", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/tree")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.AddTreeToEstate(c)) {
            assert.Equal(t, http.StatusCreated, rec.Code)
            var response models.Tree
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.NotEqual(t, uuid.Nil, response.ID)
            assert.Equal(t, tree.EstateID, response.EstateID)
            assert.Equal(t, tree.X, response.X)
            assert.Equal(t, tree.Y, response.Y)
            assert.Equal(t, tree.Height, response.Height)
        }
    })

    t.Run("invalid input format", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodPost, "/estate/"+uuid.New().String()+"/tree", bytes.NewReader([]byte("invalid")))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/tree")
        c.SetParamNames("id")
        c.SetParamValues(uuid.New().String())

        if assert.NoError(t, handler.AddTreeToEstate(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid input format", response["message"])
        }
    })

    t.Run("invalid estate ID format", func(t *testing.T) {
        tree := &models.Tree{
            X:      1,
            Y:      1,
            Height: 10,
        }
        body, _ := json.Marshal(tree)
        req := httptest.NewRequest(http.MethodPost, "/estate/invalid-id/tree", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/tree")
        c.SetParamNames("id")
        c.SetParamValues("invalid-id")

        if assert.NoError(t, handler.AddTreeToEstate(c)) {
            assert.Equal(t, http.StatusBadRequest, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Invalid estate ID format", response["message"])
        }
    })

    t.Run("estate not found", func(t *testing.T) {
        estateID := uuid.New()
        tree := &models.Tree{
            EstateID: estateID,
            X:        1,
            Y:        1,
            Height:   10,
        }
        mockEstateRepo.EXPECT().GetEstateByID(estateID).Return(nil, nil)

        body, _ := json.Marshal(tree)
        req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID.String()+"/tree", bytes.NewReader(body))
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetPath("/estate/:id/tree")
        c.SetParamNames("id")
        c.SetParamValues(estateID.String())

        if assert.NoError(t, handler.AddTreeToEstate(c)) {
            assert.Equal(t, http.StatusNotFound, rec.Code)
            var response map[string]string
            assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
            assert.Equal(t, "Estate not found", response["message"])
        }
    })
}