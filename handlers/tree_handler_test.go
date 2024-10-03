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

func TestTreeHandler_AddTreeToEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 100, Length: 200}, nil)
    mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 10, 20).Return(nil, nil)
    mockTreeRepo.EXPECT().AddTreeToEstate(gomock.Any()).Return(nil)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
        var response models.Tree
        if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
            assert.NotEmpty(t, response.ID)
        }
    }
}

func TestTreeHandler_AddTreeToEstate_InvalidInput(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": -1, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_EstateNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, nil)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_TreeAlreadyExists(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 100, Length: 200}, nil)
    mockTreeRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 10, 20).Return(&models.Tree{ID: uuid.New()}, nil)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_InvalidEstateID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    invalidEstateID := "invalid-uuid"
    req := httptest.NewRequest(http.MethodPost, "/estate/"+invalidEstateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(invalidEstateID)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_DatabaseError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(nil, errors.New("database error"))

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusInternalServerError, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_TreeCoordinatesOutOfBounds(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 150, "y": 250, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    mockEstateRepo.EXPECT().GetEstateByID(gomock.Any()).Return(&models.Estate{ID: uuid.MustParse(estateID), Width: 100, Length: 200}, nil)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}

func TestTreeHandler_AddTreeToEstate_InvalidTreeHeight(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTreeRepo := mocks.NewMockTreeRepository(ctrl)
    mockEstateRepo := mocks.NewMockEstateRepository(ctrl)
    handler := NewTreeHandler(mockTreeRepo, mockEstateRepo)

    e := echo.New()
    estateID := uuid.New().String()
    req := httptest.NewRequest(http.MethodPost, "/estate/"+estateID+"/tree", strings.NewReader(`{"x": 10, "y": 20, "height": 35}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(estateID)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }
}