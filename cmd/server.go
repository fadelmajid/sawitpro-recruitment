package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "sawitpro-recruitment/generated"
    "sawitpro-recruitment/repositories"
    "github.com/google/uuid" // Import the UUID package
)

type Server struct {
    estateRepo repositories.EstateRepository
    treeRepo   repositories.TreeRepository
}

func NewServer(estateRepo repositories.EstateRepository, treeRepo repositories.TreeRepository) *Server {
    return &Server{
        estateRepo: estateRepo,
        treeRepo:   treeRepo,
    }
}

func (s *Server) PostEstate(ctx echo.Context) error {
    var estate generated.Estate
    if err := ctx.Bind(&estate); err != nil {
        errMsg := err.Error() // Assign error message to a variable
        return ctx.JSON(http.StatusBadRequest, generated.Error{Message: &errMsg})
    }

    // Add business logic to create an estate
    // Example: s.estateRepo.CreateEstate(estate)

    return ctx.JSON(http.StatusCreated, estate)
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id uuid.UUID, params generated.GetEstateIdDronePlanParams) error {
    // Add business logic to calculate the drone's total travel distance
    // Example: dronePlan := s.droneHandler.CalculateDronePlan(id, params.MaxDistance)

    var dronePlan generated.DronePlan
    return ctx.JSON(http.StatusOK, dronePlan)
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id uuid.UUID) error {
    // Add business logic to get stats of trees in an estate
    // Example: stats := s.treeRepo.GetEstateStats(id)

    var stats generated.EstateStats
    return ctx.JSON(http.StatusOK, stats)
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id uuid.UUID) error {
    var tree generated.Tree
    if err := ctx.Bind(&tree); err != nil {
        errMsg := err.Error() // Assign error message to a variable
        return ctx.JSON(http.StatusBadRequest, generated.Error{Message: &errMsg})
    }

    // Add business logic to add a tree to an estate
    // Example: s.treeRepo.AddTreeToEstate(id, tree)

    return ctx.JSON(http.StatusCreated, tree)
}