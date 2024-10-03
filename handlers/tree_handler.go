package handlers

import (
    "net/http"
    "sawitpro-recruitment/models"
    "sawitpro-recruitment/repositories"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/sirupsen/logrus"
)

// TreeHandler manages tree-related requests.
type TreeHandler struct {
    TreeRepo   repositories.TreeRepository
    EstateRepo repositories.EstateRepository
}

// NewTreeHandler creates a new TreeHandler.
func NewTreeHandler(treeRepo repositories.TreeRepository, estateRepo repositories.EstateRepository) *TreeHandler {
    return &TreeHandler{
        TreeRepo:   treeRepo,
        EstateRepo: estateRepo,
    }
}

// AddTreeToEstate adds a tree to an existing estate
// @Summary Add a tree to an estate
// @Description Add a tree to an estate
// @Tags trees
// @Accept json
// @Produce json
// @Param id path string true "Estate ID"
// @Param tree body models.Tree true "Tree"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /estate/{id}/tree [post]
func (h *TreeHandler) AddTreeToEstate(c echo.Context) error {
    estateID := c.Param("id")
    tree := new(models.Tree)
    
    // Bind the request body to the tree model
    if err := c.Bind(tree); err != nil {
        logrus.Warnf("Failed to bind tree: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid input format",
        })
    }

    // Validate tree dimensions and height
    if tree.X < 1 || tree.Y < 1 || tree.Height < 1 || tree.Height > 30 {
        logrus.Warnf("Invalid tree coordinates or height: x=%d, y=%d, height=%d", tree.X, tree.Y, tree.Height)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid tree coordinates or height",
        })
    }

    // Convert estate ID to UUID
    estateUUID, err := uuid.Parse(estateID)
    if err != nil {
        logrus.Warnf("Invalid estate ID format: %s", estateID)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid estate ID format",
        })
    }

    // Check if the estate exists
    estate, err := h.EstateRepo.GetEstateByID(estateUUID)
    if err != nil {
        logrus.Errorf("Database error while retrieving estate ID %s: %v", estateUUID, err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "message": "Database error while retrieving estate",
        })
    }
    if estate == nil {
        logrus.Warnf("Estate not found: %s", estateUUID)
        return c.JSON(http.StatusNotFound, map[string]string{
            "message": "Estate not found",
        })
    }

    // Validate coordinates within estate bounds
    if tree.X > estate.Width || tree.Y > estate.Length {
        logrus.Warnf("Tree coordinates out of bounds: x=%d, y=%d", tree.X, tree.Y)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Tree coordinates out of bounds",
        })
    }

    // Check if a tree already exists at the given coordinates
    existingTree, err := h.TreeRepo.GetTreeByCoordinates(estateUUID, tree.X, tree.Y)
    if err != nil {
        logrus.Errorf("Database error while checking existing tree for estate ID %s: %v", estateUUID, err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "message": "Database error while checking existing tree",
        })
    }
    if existingTree != nil {
        logrus.Warnf("A tree already exists at location x=%d, y=%d for estate ID %s", tree.X, tree.Y, estateUUID)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "A tree already exists at this location",
        })
    }

    // Assign a new UUID to the tree and set the estate ID
    tree.ID = uuid.New()
    tree.EstateID = estateUUID

    // Add tree to the estate via repository
    if err := h.TreeRepo.AddTreeToEstate(tree); err != nil {
        logrus.Errorf("Failed to store tree in database for estate ID %s: %v", estateUUID, err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "message": "Failed to store tree in database",
        })
    }

    logrus.Infof("Tree added successfully to estate ID %s: %v", estateUUID, tree.ID)
    return c.JSON(http.StatusCreated, map[string]string{
        "Id": tree.ID.String(),
    })
}