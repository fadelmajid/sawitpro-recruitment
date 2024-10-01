package routes

import (
	"sawitpro-recruitment/handlers"
	"github.com/labstack/echo/v4"
)

// InitRoutes initializes the API routes.
func InitRoutes(e *echo.Echo, estateHandler *handlers.EstateHandler, treeHandler *handlers.TreeHandler, droneHandler *handlers.DroneHandler) {
	e.POST("/estate", estateHandler.CreateEstate)
	e.POST("/estate/:id/tree", treeHandler.AddTreeToEstate)
	e.GET("/estate/:id/stats", estateHandler.GetEstateStats)
	e.GET("/estate/:id/drone-plan", droneHandler.CalculateDronePlanWithLimit)
}
