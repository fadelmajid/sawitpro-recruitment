package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "sawitpro-recruitment/repositories"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/sirupsen/logrus"
)

// DroneHandler manages drone-related requests.
type DroneHandler struct {
    TreeRepo repositories.TreeRepository
}

// NewDroneHandler creates a new DroneHandler.
func NewDroneHandler(treeRepo repositories.TreeRepository) *DroneHandler {
    return &DroneHandler{
        TreeRepo: treeRepo,
    }
}

// CalculateDronePlanWithLimit calculates the drone's total travel distance with an optional max_distance parameter
// @Summary Calculate the drone's total travel distance with an optional max_distance parameter
// @Description Calculate the drone's total travel distance with an optional max_distance parameter
// @Tags drones
// @Produce json
// @Param id path string true "Estate ID"
// @Param max_distance query int false "Maximum distance the drone can travel"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /estate/{id}/drone-plan [get]
func (h *DroneHandler) CalculateDronePlanWithLimit(c echo.Context) error {
    estateID := c.Param("id")
    maxDistanceStr := c.QueryParam("max_distance")

    logrus.WithFields(logrus.Fields{
        "estateID":     estateID,
        "max_distance": maxDistanceStr,
    }).Info("Received request to calculate drone plan")

    var maxDistance int
    var limitReached bool
    if maxDistanceStr != "" {
        var err error
        maxDistance, err = strconv.Atoi(maxDistanceStr)
        if err != nil || maxDistance <= 0 {
            logrus.WithFields(logrus.Fields{
                "max_distance": maxDistanceStr,
            }).Warn("Invalid max_distance value")
            return c.JSON(http.StatusBadRequest, map[string]string{
                "message": "Invalid max_distance value",
            })
        }
    }

    // Convert to UUID
    estateUUID, err := uuid.Parse(estateID)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "estateID": estateID,
        }).Warn("Invalid estate ID format")
        return c.JSON(http.StatusBadRequest, "Invalid estate ID format")
    }

    // Get tree heights from the repository
    treeHeights, err := h.TreeRepo.GetTreesByEstateID(estateUUID)
    if err != nil {
        logrus.WithFields(logrus.Fields{
            "estateID": estateID,
        }).Error("Database error while fetching tree heights")
        return c.JSON(http.StatusInternalServerError, "Database error")
    }

    logrus.WithFields(logrus.Fields{
        "estateID": estateID,
    }).Info("Fetched tree heights")

    totalDistance := 0
    prevHeight := 0 // Start at ground level
    var landingPlotX, landingPlotY int

    // Simulate drone movement (zigzag)
    for y := 1; y <= 50; y++ { // assuming estate length and width are both 50 for now
        for x := 1; x <= 50; x++ {
            key := fmt.Sprintf("%d,%d", x, y)
            treeHeight := treeHeights[key]
            verticalDistance := abs(treeHeight - prevHeight)
            horizontalDistance := 10 // Each plot is 10 meters apart horizontally
            totalDistance += horizontalDistance + verticalDistance
            prevHeight = treeHeight

            if maxDistance > 0 && totalDistance > maxDistance {
                landingPlotX, landingPlotY = x, y
                limitReached = true
                logrus.WithFields(logrus.Fields{
                    "landingPlotX": landingPlotX,
                    "landingPlotY": landingPlotY,
                    "totalDistance": totalDistance,
                }).Info("Max distance reached")
                break
            }
        }
        if limitReached {
            break
        }
    }

    if limitReached {
        logrus.WithFields(logrus.Fields{
            "landingPlotX": landingPlotX,
            "landingPlotY": landingPlotY,
            "totalDistance": totalDistance,
        }).Info("Drone landed")
        return c.JSON(http.StatusOK, map[string]interface{}{
            "total_distance": totalDistance,
            "landed_at": map[string]int{
                "x": landingPlotX,
                "y": landingPlotY,
            },
        })
    }

    logrus.WithFields(logrus.Fields{
        "totalDistance": totalDistance,
    }).Info("Drone completed the plan")
    return c.JSON(http.StatusOK, map[string]interface{}{
        "total_distance": totalDistance,
    })
}

// Helper function for absolute value
func abs(a int) int {
    if a < 0 {
        return -a
    }
    return a
}