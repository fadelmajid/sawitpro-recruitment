package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"sawitpro-recruitment/database"
	"sawitpro-recruitment/handlers"
	"sawitpro-recruitment/repositories"
	"sawitpro-recruitment/routes"
)

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize database
	database.InitDB()

	// Create tables if they don't exist
	err := database.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	estateRepo := repositories.NewEstateRepository(database.DB)
	treeRepo := repositories.NewTreeRepository(database.DB)

	// Initialize handlers
	estateHandler := handlers.NewEstateHandler(estateRepo)
	treeHandler := handlers.NewTreeHandler(treeRepo)
	droneHandler := handlers.NewDroneHandler(treeRepo)

	// Initialize routes
	routes.InitRoutes(e, estateHandler, treeHandler, droneHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
