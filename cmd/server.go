package main

import (
	"sawitpro-recruitment/generated"
	"sawitpro-recruitment/handlers"
	"sawitpro-recruitment/repositories"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Server struct {
	estateHandler *handlers.EstateHandler
	droneHandler  *handlers.DroneHandler
	treeHandler   *handlers.TreeHandler
}

func NewServer(estateRepo repositories.EstateRepository, treeRepo repositories.TreeRepository) *Server {
	return &Server{
		estateHandler: handlers.NewEstateHandler(estateRepo),
		droneHandler:  handlers.NewDroneHandler(treeRepo, estateRepo),
		treeHandler:   handlers.NewTreeHandler(treeRepo, estateRepo),
	}
}

func (s *Server) PostEstate(ctx echo.Context) error {
	return s.estateHandler.CreateEstate(ctx)
}

func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id uuid.UUID, params generated.GetEstateIdDronePlanParams) error {
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())
	if params.MaxDistance != nil {
		ctx.QueryParams().Set("max_distance", strconv.Itoa(*params.MaxDistance))
	}
	return s.droneHandler.CalculateDronePlanWithLimit(ctx)
}

func (s *Server) GetEstateIdStats(ctx echo.Context, id uuid.UUID) error {
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())
	return s.estateHandler.GetEstateStats(ctx)
}

func (s *Server) PostEstateIdTree(ctx echo.Context, id uuid.UUID) error {
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())
	return s.treeHandler.AddTreeToEstate(ctx)
}
