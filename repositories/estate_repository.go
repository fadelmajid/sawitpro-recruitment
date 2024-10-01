package repositories

import (
	"database/sql"
	"sawitpro-recruitment/models"
	"github.com/google/uuid"
)

// EstateRepository defines the methods for estate-related database operations.
type EstateRepository interface {
	CreateEstate(estate *models.Estate) error
	GetEstateByID(id uuid.UUID) (*models.Estate, error)
	GetEstateStats(id uuid.UUID) (int, int, int, int, error)
}

// estateRepository is the concrete implementation of the EstateRepository interface.
type estateRepository struct {
	db *sql.DB
}

// NewEstateRepository returns a new instance of estateRepository.
func NewEstateRepository(db *sql.DB) EstateRepository {
	return &estateRepository{
		db: db,
	}
}

// CreateEstate inserts a new estate into the database.
func (r *estateRepository) CreateEstate(estate *models.Estate) error {
	_, err := r.db.Exec("INSERT INTO estates (id, width, length) VALUES ($1, $2, $3)", estate.ID, estate.Width, estate.Length)
	return err
}

// GetEstateByID retrieves an estate by its ID.
func (r *estateRepository) GetEstateByID(id uuid.UUID) (*models.Estate, error) {
	estate := &models.Estate{}
	err := r.db.QueryRow("SELECT id, width, length FROM estates WHERE id = $1", id).Scan(&estate.ID, &estate.Width, &estate.Length)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return estate, nil
}

// GetEstateStats retrieves statistics about trees in a specified estate.
func (r *estateRepository) GetEstateStats(id uuid.UUID) (int, int, int, int, error) {
	var count, max, min, median int
	err := r.db.QueryRow("SELECT COUNT(*), MAX(height), MIN(height), PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height) FROM trees WHERE estate_id = $1", id).Scan(&count, &max, &min, &median)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return count, max, min, median, nil
}
