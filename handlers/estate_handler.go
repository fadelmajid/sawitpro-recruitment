package handlers

import (
    "net/http"
    "sawitpro-recruitment/models"
    "sawitpro-recruitment/repositories"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// EstateHandler manages estate-related requests.
type EstateHandler struct {
    EstateRepo repositories.EstateRepository
}

// NewEstateHandler creates a new EstateHandler.
func NewEstateHandler(repo repositories.EstateRepository) *EstateHandler {
    return &EstateHandler{
        EstateRepo: repo,
    }
}

// CreateEstate handles the creation of a new estate
// @Summary Create a new estate
// @Description Create a new estate
// @Tags estates
// @Accept json
// @Produce json
// @Param estate body models.Estate true "Estate"
// @Success 201 {object} models.Estate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /estate [post]
func (h *EstateHandler) CreateEstate(c echo.Context) error {
    estate := new(models.Estate)
    if err := c.Bind(estate); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Invalid input format",
        })
    }

    // Validate estate dimensions
    if estate.Width < 1 || estate.Length < 1 || estate.Width > 50000 || estate.Length > 50000 {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "message": "Estate dimensions must be between 1 and 50000",
        })
    }

    estate.ID = uuid.New()

    // Call the repository to create estate
    if err := h.EstateRepo.CreateEstate(estate); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "message": "Failed to store estate in database",
        })
    }

    return c.JSON(http.StatusCreated, estate)
}

// GetEstateStats retrieves stats of trees in an estate
// @Summary Get stats of trees in an estate
// @Description Get stats of trees in an estate
// @Tags estates
// @Produce json
// @Param id path string true "Estate ID"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /estate/{id}/stats [get]
func (h *EstateHandler) GetEstateStats(c echo.Context) error {
    id := c.Param("id")

    // Convert to UUID
    estateID, err := uuid.Parse(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid estate ID format")
    }

    // Call the repository to get stats
    count, max, min, median, err := h.EstateRepo.GetEstateStats(estateID)
    if err != nil {
		logrus.Errorf("Failed to get estate stats for ID %s: %v", estateID, err)
        return c.JSON(http.StatusInternalServerError, "Database error")
    }

    stats := map[string]int{
        "count":  count,
        "max":    max,
        "min":    min,
        "median": median,
    }

    return c.JSON(http.StatusOK, stats)
}