package handlers

import (
    "net/http"
    "sawitpro-recruitment/models"
    "sawitpro-recruitment/repositories"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
)

// TreeHandler manages tree-related requests.
type TreeHandler struct {
    TreeRepo repositories.TreeRepository
}

// NewTreeHandler creates a new TreeHandler.
func NewTreeHandler(repo repositories.TreeRepository) *TreeHandler {
    return &TreeHandler{
        TreeRepo: repo,
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
// @Success 201 {object} models.Tree
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /estate/{id}/tree [post]
func (h *TreeHandler) AddTreeToEstate(c echo.Context) error {
    estateID := c.Param("id")
    tree := new(models.Tree)
    if err := c.Bind(tree); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid input format",
        })
    }

    // Validate tree dimensions and height
    if tree.X < 1 || tree.Y < 1 || tree.Height < 1 || tree.Height > 30 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid tree coordinates or height",
        })
    }

    // Convert to UUID
    estateUUID, err := uuid.Parse(estateID)
    if err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid estate ID format")
    }

    // Check if a tree already exists at the given coordinates
    existingTree, err := h.TreeRepo.GetTreeByCoordinates(estateUUID, tree.X, tree.Y)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, "Database error")
    }
    if existingTree != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "A tree already exists at this location",
        })
    }

    tree.ID = uuid.New()
    tree.EstateID = estateUUID

    // Add tree to the estate via repository
    if err := h.TreeRepo.AddTreeToEstate(tree); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "message": "Failed to store tree in database",
        })
    }

    return c.JSON(http.StatusCreated, tree)
}