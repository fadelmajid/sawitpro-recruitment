package tests

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "sawitpro-recruitment/handlers"
    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestTreeHandler_AddTreeToEstate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockTreeRepository(ctrl)
    handler := handlers.NewTreeHandler(mockRepo)

    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/estates/:id/trees", strings.NewReader(`{"x": 10, "y": 20, "height": 15}`))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(uuid.New().String())

    mockRepo.EXPECT().GetTreeByCoordinates(gomock.Any(), 10, 20).Return(nil, nil)
    mockRepo.EXPECT().AddTreeToEstate(gomock.Any()).Return(nil)

    if assert.NoError(t, handler.AddTreeToEstate(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
    }
}