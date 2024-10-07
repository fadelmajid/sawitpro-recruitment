package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HelloHandler(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Bad Request"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello User " + id})
}
