package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sawitpro-recruitment/handlers"
	"sawitpro-recruitment/mocks"
	"sawitpro-recruitment/models"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateEstate(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEstateRepository(ctrl)
	handler := handlers.NewEstateHandler(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		estate := &models.Estate{
			Width:  100,
			Length: 200,
		}
		mockRepo.EXPECT().CreateEstate(gomock.Any()).Return(nil)

		body, _ := json.Marshal(estate)
		req := httptest.NewRequest(http.MethodPost, "/estate", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, handler.CreateEstate(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var response map[string]string
			if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
				assert.NotEmpty(t, response["Id"])
			}
		}
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
