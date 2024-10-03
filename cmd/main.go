package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "sawitpro-recruitment/database"
    "sawitpro-recruitment/generated"
    "sawitpro-recruitment/repositories"
)

// @title SawitPro Recruitment API
// @version 1.0
// @description This is the API documentation for SawitPro Recruitment project.
// @host localhost:8080
// @BasePath /
func main() {
    // Initialize Echo instance
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Initialize database
    database.InitDB()

    // Initialize repositories
    estateRepo := repositories.NewEstateRepository(database.DB)
    treeRepo := repositories.NewTreeRepository(database.DB)

    // Initialize server
    server := NewServer(estateRepo, treeRepo)

    // Register handlers
    generated.RegisterHandlers(e, server)

    // Start server
    e.Logger.Fatal(e.Start(":8080"))
}